package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kevingentile/chet/pkg/chat"
)

type Message struct {
	// Type string `json:"type"`
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type RoomRelay struct {
	RoomID string
	// RoomName     string
	Subscription *amqp.Subscriber
	Publisher    *amqp.Publisher
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

	go runMessageBroker()
	go handleRoomSubscriptions()

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

func runMessageBroker() {
	logger := watermill.NewStdLogger(false, false)
	cqrsMarshaler := cqrs.JSONMarshaler{}

	commandsAMQPConfig := amqp.NewDurableQueueConfig(amqpAddress)

	commandsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
	}
	commandsSubscriber, err := amqp.NewSubscriber(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
	}

	eventsPublisher, err := amqp.NewPublisher(amqp.NewDurablePubSubConfig(amqpAddress, nil), logger)
	if err != nil {
		panic(err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddMiddleware(middleware.Recoverer)

	// cqrs.Facade is facade for Command and Event buses and processors.
	// You can use facade, or create buses and processors manually (you can inspire with cqrs.NewFacade)
	_, err = cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				CreateRoomPublisherCmdHandler{eb},
			}
		},
		CommandsPublisher: commandsPublisher,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			// we can reuse subscriber, because all commands have separated topics
			return commandsSubscriber, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			// because we are using PubSub RabbitMQ config, we can use one topic for all events
			return "chet-events"
		},
		// EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
		// 	return []cqrs.EventHandler{}
		// },
		EventsPublisher: eventsPublisher,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			config := amqp.NewDurablePubSubConfig(
				amqpAddress,
				amqp.GenerateQueueNameTopicNameWithSuffix(handlerName),
			)

			return amqp.NewSubscriber(config, logger)
		},
		Router:                router,
		CommandEventMarshaler: cqrsMarshaler,
		Logger:                logger,
	})
	if err != nil {
		panic(err)
	}

	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}

}

type CreateRoomPublisherCmdHandler struct {
	eventBus *cqrs.EventBus
}

func (h CreateRoomPublisherCmdHandler) HandlerName() string {
	return "CreateRoomPublisherCmdHandler"
}

func (h CreateRoomPublisherCmdHandler) NewCommand() interface{} {
	return &chat.CreateRoomPublisherCmd{}
}

func (h CreateRoomPublisherCmdHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*chat.CreateRoomPublisherCmd)
	log.Println("Creating room subscriber")
	if err := addRoomSubscriber(cmd.ID); err != nil {
		return err
	}

	// if err := h.eventBus.Publish(ctx, &chat.RoomCreatedEvent{
	// 	Room: *room,
	// 	Time: time.Now().UTC(),
	// }); err != nil {
	// 	return err
	// }

	// log.Println("Created chat room:", *room)
	return nil
}

func addRoomSubscriber(roomID chat.RoomID) error {
	for _, r := range roomRelays {
		if r.RoomID == roomID.String() {
			log.Println("Room Subscriber already exists")
			return nil
		}
	}
	cfg := amqp.NewNonDurablePubSubConfig(amqpAddress, func(n string) string { return roomID.String() })
	subscriber, err := amqp.NewSubscriber(cfg, nil)
	if err != nil {
		return err
	}

	publisher, err := amqp.NewPublisher(cfg, nil)
	if err != nil {
		return err
	}

	relay := &RoomRelay{
		RoomID:       roomID.String(),
		Subscription: subscriber,
		Publisher:    publisher,
	}

	roomRelays = append(roomRelays, relay)
	log.Println("Added room relay", relay.RoomID)

	//TODO move to chat connect
	// msgChanel, err := subscriber.Subscribe(context.Background(), roomID.String())
	// if err != nil {
	// 	return err
	// }
	// messageChannels = append(messageChannels, msgChanel)

	// testMsg := &Message{
	// 	ID:        "asdf",
	// 	Timestamp: time.Now().Unix(),
	// }
	// payload, err := json.Marshal(testMsg)
	// if err != nil {
	// 	return err
	// }
	// message := message.NewMessage(watermill.NewUUID(), payload)
	// publisher.Publish(queueName, message)

	return nil
}

func handleRoomSubscriptions() {
	for _, c := range messageChannels {
		msg := <-c
		msg.Ack()

		message := &Message{}
		if err := json.Unmarshal(msg.Payload, message); err != nil {
			log.Println(err)
		}

		fmt.Println("Received room message: ", message)
	}
}
