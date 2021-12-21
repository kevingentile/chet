package eventstream

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type EventHandler interface {
	EventName() string
	Handle(ctx context.Context, data []byte) error
}

type EventStream struct {
	handlers      []EventHandler
	messageStream *MessageStream
	streamContext context.Context
	producer      *MessageStreamProducer
	consumer      *MessageStreamConsumer
	offsetSpec    stream.OffsetSpecification
}

func NewEventStream(ms *MessageStream, offsetSpec stream.OffsetSpecification) (*EventStream, error) {
	es := &EventStream{
		handlers:      make([]EventHandler, 0),
		messageStream: ms,
		offsetSpec:    offsetSpec,
	}
	return es, nil
}

func (es *EventStream) AddEventHandler(h EventHandler) {
	es.handlers = append(es.handlers, h)
}

func (es *EventStream) HandlerNames() []string {
	names := make([]string, len(es.handlers))
	for i, h := range es.handlers {
		names[i] = h.EventName()
	}
	return names
}

func (es *EventStream) GetEventHandler(name string) EventHandler {
	for _, h := range es.handlers {
		if h.EventName() == name {
			return h
		}
	}
	return nil
}

func (es *EventStream) Run(ctx context.Context) error {
	es.streamContext = ctx

	producer, err := es.messageStream.NewProducer(&MessageStreamProducerConfig{})
	if err != nil {
		return err
	}

	es.producer = producer

	consumer, err := es.messageStream.NewConsumer(&MessageStreamConsumerConfig{
		Options: stream.NewConsumerOptions().SetAutoCommit(stream.NewAutoCommitStrategy().
			SetCountBeforeStorage(50). // store each 50 messages stores
			SetFlushInterval(10 * time.Second)),
		MessageHandler: es.messageHandler,
		Offset:         stream.OffsetSpecification{}.First(),
	})
	if err != nil {
		return err
	}

	es.consumer = consumer
	<-ctx.Done()
	return err
}

func (es *EventStream) Publish(subject string, message *amqp.AMQP10) error {
	message.Properties = &amqp.MessageProperties{Subject: subject}
	if err := es.producer.Send(message); err != nil {
		return err
	}

	return nil
}

func (es *EventStream) messageHandler(consumerContext stream.ConsumerContext, msg *amqp.Message) {
	handler := es.GetEventHandler(msg.Properties.Subject)
	if handler == nil {
		log.Println("skipping subject", msg.Properties.Subject)
		return
	}

	for _, d := range msg.Data {
		ctx, cancelFunc := context.WithTimeout(es.streamContext, time.Duration(time.Minute*1))
		if err := handler.Handle(ctx, d); err != nil {
			log.Println(err)
		}
		cancelFunc()
	}

	// consumerContext.Consumer.StoreOffset()
}
