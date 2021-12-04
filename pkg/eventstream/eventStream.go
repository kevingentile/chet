package eventstream

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type EventStream struct {
	Handlers      []EventHandler
	MessageStream *MessageStream
	streamContext context.Context
	producer      *MessageStreamProducer
	consumer      *MessageStreamConsumer
}

func NewEventStream(stream *MessageStream) (*EventStream, error) {
	es := &EventStream{
		Handlers:      make([]EventHandler, 0),
		MessageStream: stream,
	}
	return es, nil
}

func (es *EventStream) AddEventHandler(h EventHandler) {
	es.Handlers = append(es.Handlers, h)
}

func (es *EventStream) HandlerNames() []string {
	names := make([]string, len(es.Handlers))
	for i, h := range es.Handlers {
		names[i] = h.EventName()
	}
	return names
}

func (es *EventStream) GetEventHandler(name string) EventHandler {
	for _, h := range es.Handlers {
		if h.EventName() == name {
			return h
		}
	}
	return nil
}

func (es *EventStream) Run(ctx context.Context) error {
	es.streamContext = ctx

	producer, err := es.MessageStream.NewProducer(&MessageStreamProducerConfig{})
	if err != nil {
		return err
	}
	es.producer = producer

	consumer, err := es.MessageStream.NewConsumer(&MessageStreamConsumerConfig{
		Options:        stream.NewConsumerOptions().SetOffset(stream.OffsetSpecification{}.First()),
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

func (es *EventStream) messageHandler(consumerContext stream.ConsumerContext, msg *amqp.Message) {
	handler := es.GetEventHandler(msg.Properties.Subject)
	if handler == nil {
		log.Println("skipping subject", msg.Properties.Subject)
		return
	}

	for _, d := range msg.Data {
		ctx, _ := context.WithTimeout(es.streamContext, time.Duration(time.Minute*1))
		if err := handler.Handle(ctx, d); err != nil {
			log.Println(err)
		}
	}

	// consumerContext.Consumer.StoreOffset()
}
