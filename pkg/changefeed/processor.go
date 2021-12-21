package changefeed

import (
	"context"

	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type Processor struct {
	monitorMessage *eventstream.MessageStream
	monitorEvent   *eventstream.EventStream

	leaseMessage *eventstream.MessageStream
	leaseEvent   *eventstream.EventStream
}

type Config struct {
	Name                 string
	MonitorEnvOptions    *stream.EnvironmentOptions
	MonitorStreamOptions *stream.StreamOptions

	LeaseEnvOptions    *stream.EnvironmentOptions
	LeaseStreamOptions *stream.StreamOptions
}

func NewProcessor(config *Config) (*Processor, error) {
	leaseStream, err := eventstream.NewMessageStream(&eventstream.StreamConfig{
		EnvOptions:    config.LeaseEnvOptions,
		StreamOptions: config.LeaseStreamOptions,

		StreamName: config.Name + "-lease",
	})
	if err != nil {
		return nil, err
	}

	monitorStream, err := eventstream.NewMessageStream(&eventstream.StreamConfig{
		EnvOptions:    config.MonitorEnvOptions,
		StreamOptions: config.MonitorStreamOptions,
		StreamName:    config.Name,
	})
	if err != nil {
		return nil, err
	}

	return &Processor{
		monitorMessage: monitorStream,
		leaseMessage:   leaseStream,
	}, nil
}

func (p *Processor) Run(ctx context.Context) error {
	leaseStream, err := eventstream.NewEventStream(p.leaseMessage, stream.OffsetSpecification{}.Last())
	if err != nil {
		return err
	}

	return nil
}
