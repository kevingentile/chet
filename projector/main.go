package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kevingentile/chet/pkg/chat"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/kevingentile/chet/pkg/infrastructure"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoUri = "mongodb://root:example@localhost:27017/?maxPoolSize=20&w=majority"

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	messageStream, err := eventstream.NewMessageStream(&eventstream.StreamConfig{
		EnvOptions: stream.NewEnvironmentOptions().
			SetHost("localhost").
			SetPort(5552).
			SetUser("guest").
			SetPassword("guest"),
		StreamOptions: &stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.MB(500),
		},
		StreamName: "chet-messages",
	})
	if err != nil {
		panic(err)
	}

	eventStream, err := eventstream.NewEventStream(messageStream)
	if err != nil {
		panic(err)
	}

	roomStore, err := infrastructure.NewDefaultMongoStore(client, "chet", "room")
	if err != nil {
		panic(err)
	}

	eventStream.AddEventHandler(&roomCreatedHandler{
		store: roomStore,
	})

	if err := eventStream.Run(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

type roomCreatedHandler struct {
	store *infrastructure.MongoStore
}

func (h *roomCreatedHandler) EventName() string {
	return chat.RoomCreatedEvent
}

func (h *roomCreatedHandler) Handle(ctx context.Context, data []byte) error {
	e := &chat.RoomCreated{}
	if err := json.Unmarshal(data, e); err != nil {
		return err
	}

	view := chat.NewRoomView(e.Room)
	if err := h.store.Create(view); err != nil {
		return err
	}

	log.Println("Created chat room view:", view)

	return nil
}
