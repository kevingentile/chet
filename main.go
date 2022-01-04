package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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
			time.Sleep(20 * time.Second)
		}
	}()

	leaseID := "main-" + uuid.NewString()
	for {
		// check for an existing lease
		log.Println("checking lease...")
		lastLease := &Lease{}
		auxCollection := database.Collection("chet-aux")
		opt := options.FindOne()
		opt.SetSort(bson.M{"$natural": -1})
		if err := auxCollection.FindOne(context.TODO(), bson.M{}, opt).Decode(lastLease); err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				panic(err)
			}
		}

		// create a new lease
		lease := Lease{
			ID:        leaseID,
			ExpiresAt: time.Now().Add(time.Minute * 1),
			StartAt:   time.Now(),
		}
		// if no lease exists, acquire the first lease
		// TODO service startup race may insert multiples
		if *lastLease == (Lease{}) {
			log.Println("acquiring initial lease..", lease)
			if _, err := auxCollection.InsertOne(context.TODO(), lease); err != nil {
				panic(err)
			}
		} else {
			// if the previous lease has expired, attempt to acquire a new lease
			if lastLease.ExpiresAt.Before(time.Now()) {
				filter := bson.M{"id": lastLease.ID}
				log.Println("acquiring lease...", lease)
				r, err := auxCollection.ReplaceOne(context.TODO(), filter, lease)
				if err != nil {
					panic(err)
				}
				if r.MatchedCount > 1 {
					panic(errors.New("ERROR: lease replacement matches multiple results"))
				}
			}
		}

		// verify
		checkLease := &Lease{}
		if err := auxCollection.FindOne(context.TODO(), bson.M{}, opt).Decode(checkLease); err != nil {
			panic(err)
		}

		if checkLease.ID != lease.ID {
			log.Println("lease not owned", *checkLease)
		} else {
			log.Println("acquired lease", checkLease.ID)
			copts := options.ChangeStream()
			copts.SetStartAtOperationTime(&primitive.Timestamp{T: uint32(checkLease.StartAt.Unix())})
			changeStream, err := streamCollection.Watch(context.TODO(), mongo.Pipeline{}, copts)
			if err != nil {
				panic(err)
			}

			runner := func() error {
				if changeStream.Next(context.TODO()) {
					var data ChangeEvent

					if err := changeStream.Decode(&data); err != nil {
						return err
					}
					log.Println("data:", data.ID, data.FullDocument)
				}
				return nil
			}
			runUntil(checkLease.ExpiresAt, runner)
			changeStream.Close(context.TODO())
		}
		time.Sleep(time.Until(checkLease.ExpiresAt))
	}
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
