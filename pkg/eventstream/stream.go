package eventstream

import (
	"errors"
	"log"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/message"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type StreamConfig struct {
	EnvOptions    *stream.EnvironmentOptions
	StreamOptions *stream.StreamOptions
	StreamName    string
}

type MessageStream struct {
	Env        *stream.Environment
	StreamName string
}

func NewMessageStream(config *StreamConfig) (*MessageStream, error) {
	if config.StreamName == "" {
		return nil, errors.New("invalid stream name")
	}

	if config.EnvOptions == nil {
		return nil, errors.New("invalid environment options")
	}

	if config.StreamOptions == nil {
		return nil, errors.New("invalid stream options")
	}

	env, err := stream.NewEnvironment(config.EnvOptions)
	if err != nil {
		return nil, err
	}

	exists, err := env.StreamExists(config.StreamName)
	if err != nil {
		return nil, err
	}

	if exists {
		log.Println("Using existing stream", config.StreamName)
	} else {
		log.Println("Initializing new stream", config.StreamName)
		if err := env.DeclareStream(config.StreamName, config.StreamOptions); err != nil {
			return nil, err
		}
	}
	return &MessageStream{
		Env:        env,
		StreamName: config.StreamName,
	}, nil
}

func (ms *MessageStream) NewProducer(config *MessageStreamProducerConfig) (*MessageStreamProducer, error) {
	producer, err := ms.Env.NewProducer(ms.StreamName, config.Options)
	if err != nil {
		return nil, err
	}

	return &MessageStreamProducer{
		Producer: producer,
	}, nil
}

func (ms *MessageStream) NewConsumer(config *MessageStreamConsumerConfig) (*MessageStreamConsumer, error) {
	consumer, err := ms.Env.NewConsumer(ms.StreamName, config.MessageHandler, config.Options)
	if err != nil {
		return nil, err
	}
	return &MessageStreamConsumer{
		Consumer: consumer,
	}, nil
}

func (ms *MessageStream) Delete() error {
	return ms.Env.DeleteStream(ms.StreamName)
}

type MessageStreamConsumerConfig struct {
	Options        *stream.ConsumerOptions
	MessageHandler stream.MessagesHandler
	Offset         stream.OffsetSpecification
}

type MessageStreamConsumer struct {
	Consumer *stream.Consumer
}

func (msc *MessageStreamConsumer) Close() error {
	return msc.Consumer.Close()
}

type MessageStreamProducerConfig struct {
	Options *stream.ProducerOptions
}

type MessageStreamProducer struct {
	Producer *stream.Producer
}

func (msp *MessageStreamProducer) Close() error {
	return msp.Producer.Close()
}

func (msp *MessageStreamProducer) Send(message message.StreamMessage) error {
	return msp.Producer.Send(message)
}