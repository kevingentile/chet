package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kevingentile/chet/pkg/chat"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/kevingentile/chet/pkg/service"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func main() {
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
		log.Fatal(err)
	}

	eventStream, err := eventstream.NewEventStream(messageStream)
	if err != nil {
		log.Fatal(err)
	}

	roomService, err := service.NewRoomService(messageStream)
	if err != nil {
		log.Fatal(err)
	}

	roomReactor, err := service.NewRoomReactor(messageStream)

	createRoomHandler := &createRoomHandler{
		roomService: roomService,
		roomReactor: roomReactor,
	}
	eventStream.AddEventHandler(createRoomHandler)
	if err := eventStream.Run(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

type createRoomHandler struct {
	roomService *service.RoomService
	roomReactor *service.RoomReactor
}

func (h *createRoomHandler) EventName() string {
	return chat.CreateRoomCmd
}

func (h *createRoomHandler) Handle(ctx context.Context, data []byte) error {
	log.Println("Handling RoomCreated")
	command := &chat.CreateRoom{}
	if err := json.Unmarshal(data, command); err != nil {
		return err
	}

	room := chat.NewRoom(command.Host)

	if err := h.roomReactor.RoomCreated(room.ID); err != nil {
		return err
	}

	return nil
}
