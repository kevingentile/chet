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

var amqpAddress = "amqp://guest:guest@localhost:5672/"

type ChatReactor struct {
	Reactor     *eventstream.Reactor
	RoomService *service.RoomService
}

func main() {
	reactor := NewChatReactor()
	if err := reactor.Reactor.Run(context.TODO()); err != nil {
		panic(err)
	}
}

func (cr *ChatReactor) HandleCreateRoom(command *chat.CreateRoom) {
	log.Println("Handling CreateRoom")
	room := chat.NewRoom(command.Host)
	event := &chat.RoomCreated{
		Room: *room,
		Time: time.Now().UTC(),
	}
	payload, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	msg := amqp.NewMessage(payload)
	msg.Properties = &amqp.MessageProperties{Subject: chat.RoomCreatedEvent}
	if err := cr.Reactor.Producer.Send(msg); err != nil {
		panic(err)
	}
}

func NewChatReactor() *ChatReactor {
	var chatReactor *ChatReactor
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

	messageProducer, err := messageStream.NewProducer(&eventstream.MessageStreamProducerConfig{})
	if err != nil {
		panic(err)
	}

	roomService, err := service.NewRoomService(messageStream)
	if err != nil {
		panic(err)
	}

	handleMessages := func(consumerContext stream.ConsumerContext, msg *amqp.Message) {
		for _, b := range msg.Data {
			switch msg.Properties.Subject {
			case chat.CreateRoomCmd:
				cr := &chat.CreateRoom{}
				if err := json.Unmarshal(b, cr); err != nil {
					panic(err)
				}
				chatReactor.HandleCreateRoom(cr)
			default:
				// panic("unmached message")
			}
		}
		consumerContext.Consumer.StoreOffset()
	}

	messageConsumer, err := messageStream.NewConsumer(&eventstream.MessageStreamConsumerConfig{
		Options: stream.NewConsumerOptions().
			SetConsumerName("chet-chat-reactor").SetOffset(stream.OffsetSpecification{}.First()),
		MessageHandler: handleMessages,
		Offset:         stream.OffsetSpecification{}.First(),
	})
	if err != nil {
		panic(err)
	}
	chatReactor = &ChatReactor{
		Reactor: &eventstream.Reactor{
			Stream:   messageStream,
			Producer: messageProducer,
			Consumer: messageConsumer,
		},
		RoomService: roomService,
	}
	return chatReactor
}

// type CreateRoomCmdHandler struct {
// 	eventBus   *cqrs.EventBus
// 	roomSerice *service.RoomService
// }

// func (h CreateRoomCmdHandler) HandlerName() string {
// 	return "CreateRoomHandler"
// }

// func (h CreateRoomCmdHandler) NewCommand() interface{} {
// 	return &chat.CreateRoomCmd{}
// }

// func (h CreateRoomCmdHandler) Handle(ctx context.Context, c interface{}) error {
// 	cmd := c.(*chat.CreateRoomCmd)

// 	room := chat.NewRoom(cmd.Host)

// 	if err := h.eventBus.Publish(ctx, &chat.RoomCreatedEvent{
// 		Room: *room,
// 		Time: time.Now().UTC(),
// 	}); err != nil {
// 		return err
// 	}

// 	log.Println("Created chat room:", *room)

// 	if err := h.roomSerice.CreateRoomPublisher(room.ID); err != nil {
// 		return err
// 	}

// 	log.Println("Initiaing create room publisher", room.ID)

// 	return nil
// }
