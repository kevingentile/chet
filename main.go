package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageType = string

const (
	MessageCreate MessageType = "CreateRoom"
)

type Message struct {
	Type    MessageType `json:"type" bson:"type"`
	Payload []byte      `json:"payload" bson:"payload"`
}

type namespace struct {
	Db   string `bson:"db"`
	Coll string `bson:"coll"`
}

type ChangeEvent struct {
	ID            primitive.D         `bson:"_id"`
	OperationType string              `bson:"operationType"`
	ClusterTime   primitive.Timestamp `bson:"clusterTime"`
	FullDocument  Message             `bson:"fullDocument"`
	// DocumentKey   primitive.ObjectID  `bson:"documentKey"`
	// Ns            namespace           `bson:"ns"`
}

type Lease struct {
	ID        string    `json:"id" bson:"id"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
	StartAt   time.Time `json:"startAt" bson:"startAt"`
}

type CreateRoomMessage struct {
	ID string `json:"roomId"`
}

const mongoUri = "mongodb://root:example@localhost:27017/?maxPoolSize=20&w=majority&connect=direct"

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	defer client.Disconnect(context.TODO())
	database := client.Database("chet")
	streamCollection := database.Collection("chet-stream")

	changeStream, err := streamCollection.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		panic(err)
	}

	defer changeStream.Close(context.TODO())
	runner := func() error {
		if changeStream.Next(context.TODO()) {
			var data ChangeEvent

			if err := changeStream.Decode(&data); err != nil {
				return err
			}
			fmt.Println("data:", data.FullDocument)
		}
		return nil
	}

	deadline := time.Now().Add(time.Second * 10)

	go func() {
		for {
			msg := Message{
				Type:    MessageCreate,
				Payload: []byte("asdf"),
			}
			fmt.Println("inserting", msg)
			m, err := bson.Marshal(msg)
			if err != nil {
				log.Println(err)
			}
			streamCollection.InsertOne(context.TODO(), m)
			time.Sleep(1 * time.Second)
		}
	}()
	runUntil(deadline, runner)
}

func runUntil(deadline time.Time, runner func() error) {
	now := time.Now()
	errors := 0
	for !now.After(deadline) && errors < 3 {
		if err := runner(); err != nil {
			log.Println(err)
			errors++
		}
		now = time.Now()
	}
}

// func main() {

// 	messageStream, err := eventstream.NewMessageStream(&eventstream.StreamConfig{
// 		EnvOptions: stream.NewEnvironmentOptions().
// 			SetHost("localhost").
// 			SetPort(5552).
// 			SetUser("guest").
// 			SetPassword("guest"),
// 		StreamOptions: &stream.StreamOptions{
// 			MaxLengthBytes: stream.ByteCapacity{}.MB(500),
// 		},
// 		StreamName: "chet-messages",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := messageStream.Reset(); err != nil {
// 		panic(err)
// 	}

// 	roomService, _ := service.NewRoomService(messageStream)

// 	for {
// 		fmt.Println("creating room")
// 		if err := roomService.CreateRoom(uuid.New()); err != nil {
// 			panic(err)
// 		}
// 		time.Sleep(time.Second * 5)
// 	}

// }
