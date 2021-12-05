package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/kevingentile/chet/pkg/chat"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/kevingentile/chet/pkg/service"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
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

	createRoomHandler := &createRoomHandler{
		roomService: roomService,
		eventStream: eventStream,
	}
	eventStream.AddEventHandler(createRoomHandler)
	if err := eventStream.Run(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

type createRoomHandler struct {
	roomService *service.RoomService
	eventStream *eventstream.EventStream
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
	event := &chat.RoomCreated{
		Room: *room,
		Time: time.Now().UTC(),
	}

	//TODO move to service?
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := amqp.NewMessage(payload)
	if err := h.eventStream.Publish(chat.RoomCreatedEvent, msg); err != nil {
		return err
	}

	return nil
}
