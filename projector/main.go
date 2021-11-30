package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kevingentile/chet/pkg/chat"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/kevingentile/chet/pkg/infrastructure"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const amqpAddress = "amqp://guest:guest@localhost:5672/"
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

	roomStore, err := infrastructure.NewDefaultMongoStore(client, "chet", "room")
	if err != nil {
		panic(err)
	}

	handleMessages := func(consumerContext stream.ConsumerContext, msg *amqp.Message) {
		for _, b := range msg.Data {
			switch msg.Properties.Subject {
			case chat.RoomCreatedEvent:
				crm := &chat.RoomCreated{}
				if err := json.Unmarshal(b, crm); err != nil {
					panic(err)
				}
				handleRoomCreated(roomStore, crm)
			default:
				// panic("unmached message")
			}
		}
		consumerContext.Consumer.StoreOffset()
	}

	messageConsumer, err := messageStream.NewConsumer(&eventstream.MessageStreamConsumerConfig{
		Options: stream.NewConsumerOptions().
			SetConsumerName("my_consumer").SetOffset(stream.OffsetSpecification{}.First()),
		MessageHandler: handleMessages,
		Offset:         stream.OffsetSpecification{}.First(),
	})
	if err != nil {
		panic(err)
	}

	projector := &eventstream.Projector{
		Stream:   messageStream,
		Consumer: messageConsumer,
	}

	projector.Run(context.TODO())
}

func handleRoomCreated(roomStore *infrastructure.MongoStore, event *chat.RoomCreated) {
	view := chat.NewRoomView(event.Room)
	if err := roomStore.Create(view); err != nil {
		log.Println(err)
		return
	}

	log.Println("Created chat room view:", view)
}

// logger := watermill.NewStdLogger(false, false)
// cqrsMarshaler := cqrs.JSONMarshaler{}

// commandsAMQPConfig := amqp.NewDurableQueueConfig(amqpAddress)
// commandsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, logger)
// if err != nil {
// 	panic(err)
// }
// commandsSubscriber, err := amqp.NewSubscriber(commandsAMQPConfig, logger)
// if err != nil {
// 	panic(err)
// }

// eventsPublisher, err := amqp.NewPublisher(amqp.NewDurablePubSubConfig(amqpAddress, nil), logger)
// if err != nil {
// 	panic(err)
// }

// router, err := message.NewRouter(message.RouterConfig{}, logger)
// if err != nil {
// 	panic(err)
// }

// router.AddMiddleware(middleware.Recoverer)

// // cqrs.Facade is facade for Command and Event buses and processors.
// // You can use facade, or create buses and processors manually (you can inspire with cqrs.NewFacade)
// _, err = cqrs.NewFacade(cqrs.FacadeConfig{
// 	GenerateCommandsTopic: func(commandName string) string {
// 		// we are using queue RabbitMQ config, so we need to have topic per command type
// 		return commandName
// 	},
// 	// CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
// 	// 	return []cqrs.CommandHandler{
// 	// 		RoomCreatedHandler{eb, roomStore},
// 	// 	}
// 	// },
// 	CommandsPublisher: commandsPublisher,
// 	CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
// 		// we can reuse subscriber, because all commands have separated topics
// 		return commandsSubscriber, nil
// 	},
// 	GenerateEventsTopic: func(eventName string) string {
// 		// because we are using PubSub RabbitMQ config, we can use one topic for all events
// 		return "chet-events"
// 	},
// 	EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
// 		return []cqrs.EventHandler{
// 			RoomCreatedHandler{eb, roomStore},
// 		}
// 	},
// 	EventsPublisher: eventsPublisher,
// 	EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
// 		config := amqp.NewDurablePubSubConfig(
// 			amqpAddress,
// 			amqp.GenerateQueueNameTopicNameWithSuffix(handlerName),
// 		)

// 		return amqp.NewSubscriber(config, logger)
// 	},
// 	Router:                router,
// 	CommandEventMarshaler: cqrsMarshaler,
// 	Logger:                logger,
// })
// if err != nil {
// 	panic(err)
// }

// if err := router.Run(context.Background()); err != nil {
// 	panic(err)
// }

// type RoomCreatedHandler struct {
// 	eventBus *cqrs.EventBus
// 	store    infrastructure.ViewStorer
// }

// func (h RoomCreatedHandler) HandlerName() string {
// 	return "RoomCreatedHandler"
// }

// func (h RoomCreatedHandler) NewEvent() interface{} {
// 	return &chat.RoomCreatedEvent{}
// }

// func (h RoomCreatedHandler) Handle(ctx context.Context, c interface{}) error {
// 	event := c.(*chat.RoomCreatedEvent)
// 	view := chat.NewRoomView(event.Room)
// 	if err := h.store.Create(view); err != nil {
// 		return err
// 	}

// 	log.Println("Created chat room view:", view)
// 	return nil
// }
