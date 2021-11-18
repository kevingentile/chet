// package main

// import (
// 	"time"

// 	"github.com/ThreeDotsLabs/watermill"
// 	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
// 	"github.com/ThreeDotsLabs/watermill/message"
// 	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
// 	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
// )

// var (
// 	brokers      = []string{"127.0.0.1:9092"}
// 	cosumerTopic = "room-events"
// )

// func main() {
// 	logger := watermill.NewStdLogger(false, false)
// 	logger.Info("Starting the consumer", nil)

// 	r, err := message.NewRouter(message.RouterConfig{}, logger)
// 	if err != nil {
// 		panic(err)
// 	}

// 	retryMiddleware := middleware.Retry{
// 		MaxRetries:      1,
// 		InitialInterval: time.Millisecond * 10,
// 	}

// 	r.AddMiddleware(
// 		// Recoverer middleware recovers panic from handlers and middlewares
// 		middleware.Recoverer,

// 		// Limit incoming messages to 10 per second
// 		middleware.NewThrottle(10, time.Second).Middleware,

// 		// Retry middleware retries message processing if an error occurred in the handler
// 		retryMiddleware.Middleware,

// 		// Correlation ID middleware adds the correlation ID of the consumed message to each produced message.
// 		// It's useful for debugging.
// 		middleware.CorrelationID,

// 		// Simulate errors or panics from handler
// 		middleware.RandomFail(0.01),
// 		middleware.RandomPanic(0.01),
// 	)

// 	// Close the router when a SIGTERM is sent
// 	r.AddPlugin(plugin.SignalsHandler)

// 	// Handler that counts consumed posts
// 	r.AddNoPublisherHandler(
// 		"chat_broker",
// 		cosumerTopic,
// 		createSubscriber(cosumerTopic, logger),
// 		messageHandler,
// 	)
// }

// func createSubscriber(consumerGroup string, logger watermill.LoggerAdapter) message.Subscriber {
// 	sub, err := kafka.NewSubscriber(
// 		kafka.SubscriberConfig{
// 			Brokers:       brokers,
// 			Unmarshaler:   kafka.DefaultMarshaler{},
// 			ConsumerGroup: consumerGroup,
// 		},
// 		logger,
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return sub
// }

// func messageHandler(msg *message.Message) error {
// 	switch msg.Payload {
// 		case
// 	}
// }
