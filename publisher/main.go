package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kevingentile/chet/pkg/chat"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/kevingentile/chet/pkg/service"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type Message struct {
	// Type string `json:"type"`
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type RoomRelay struct {
	RoomID string
	// RoomName     string
	Subscription *eventstream.MessageStreamConsumer
	Publisher    *eventstream.MessageStreamProducer
}

const amqpAddress = "amqp://guest:guest@localhost:5672/"

var roomRelays []*RoomRelay = make([]*RoomRelay, 0)
var messageChannels []<-chan *message.Message = []<-chan *message.Message{}

func main() {
	router := gin.Default()
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return r.Host == "localhost:8000"
	}

	go runMessageBroker(context.TODO())
	// go handleRoomSubscriptions()

	router.GET("/ws", func(c *gin.Context) {
		roomID := c.Query("roomId")
		if roomID == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		var relay *RoomRelay
		for _, r := range roomRelays {
			if r.RoomID == roomID {
				relay = r
				break
			}
		}
		if relay == nil {
			c.Status(http.StatusNotFound)
			return
		}
		con, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer con.Close()

		for {
			m, err := json.Marshal(Message{
				ID:        roomID,
				Timestamp: time.Now().Unix(),
			})
			if err != nil {
				log.Println(err)
			}
			if err := con.WriteMessage(websocket.TextMessage, m); err != nil {
				break
			}
			time.Sleep(2 * time.Second)
		}
	})
	router.Run("localhost:8000")
}

func runMessageBroker(ctx context.Context) {
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

	roomService, err := service.NewRoomService(messageStream)
	if err != nil {
		panic(err)
	}

	messageHandler := func(consumerContext stream.ConsumerContext, msg *amqp.Message) {
		for _, b := range msg.Data {
			switch msg.Properties.Subject {
			case chat.RoomCreatedEvent:
				crm := &chat.RoomCreated{}
				if err := json.Unmarshal(b, crm); err != nil {
					panic(err)
				}
				handleRoomCreated(roomService, crm)
			default:
			}
		}
	}

	messageConsumer, err := messageStream.NewConsumer(&eventstream.MessageStreamConsumerConfig{
		Options: stream.NewConsumerOptions().
			SetConsumerName("chet-publisher").SetOffset(stream.OffsetSpecification{}.First()),
		MessageHandler: messageHandler,
		Offset:         stream.OffsetSpecification{}.First(),
	})
	if err != nil {
		panic(err)
	}

	messageProducer, err := messageStream.NewProducer(&eventstream.MessageStreamProducerConfig{})
	if err != nil {
		panic(err)
	}

	reactor := &eventstream.Reactor{
		Stream:   messageStream,
		Producer: messageProducer,
		Consumer: messageConsumer,
	}
	if err := reactor.Run(ctx); err != nil {
		panic(err)
	}
}

func handleRoomCreated(roomService *service.RoomService, command *chat.RoomCreated) {
	fmt.Println("Handling Room Created", command)
	//TODO create subs
}

// type CreateRoomPublisherCmdHandler struct {
// 	eventBus *cqrs.EventBus
// }

// func (h CreateRoomPublisherCmdHandler) HandlerName() string {
// 	return "CreateRoomPublisherCmdHandler"
// }

// func (h CreateRoomPublisherCmdHandler) NewCommand() interface{} {
// 	return &chat.CreateRoomPublisherCmd{}
// }

// func (h CreateRoomPublisherCmdHandler) Handle(ctx context.Context, c interface{}) error {
// 	cmd := c.(*chat.CreateRoomPublisherCmd)
// 	log.Println("Creating room subscriber")
// 	if err := addRoomSubscriber(cmd.ID); err != nil {
// 		return err
// 	}

// 	// if err := h.eventBus.Publish(ctx, &chat.RoomCreatedEvent{
// 	// 	Room: *room,
// 	// 	Time: time.Now().UTC(),
// 	// }); err != nil {
// 	// 	return err
// 	// }

// 	// log.Println("Created chat room:", *room)
// 	return nil
// }

// func addRoomSubscriber(roomID chat.RoomID) error {
// 	for _, r := range roomRelays {
// 		if r.RoomID == roomID.String() {
// 			log.Println("Room Subscriber already exists")
// 			return nil
// 		}
// 	}
// 	cfg := amqp.NewNonDurablePubSubConfig(amqpAddress, func(n string) string { return roomID.String() })
// 	subscriber, err := amqp.NewSubscriber(cfg, nil)
// 	if err != nil {
// 		return err
// 	}

// 	publisher, err := amqp.NewPublisher(cfg, nil)
// 	if err != nil {
// 		return err
// 	}

// 	relay := &RoomRelay{
// 		RoomID:       roomID.String(),
// 		Subscription: subscriber,
// 		Publisher:    publisher,
// 	}

// 	roomRelays = append(roomRelays, relay)
// 	log.Println("Added room relay", relay.RoomID)

// 	//TODO move to chat connect
// 	// msgChanel, err := subscriber.Subscribe(context.Background(), roomID.String())
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// messageChannels = append(messageChannels, msgChanel)

// 	// testMsg := &Message{
// 	// 	ID:        "asdf",
// 	// 	Timestamp: time.Now().Unix(),
// 	// }
// 	// payload, err := json.Marshal(testMsg)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// message := message.NewMessage(watermill.NewUUID(), payload)
// 	// publisher.Publish(queueName, message)

// 	return nil
// }
