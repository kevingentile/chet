package eventstream

import "context"

type Reactor struct {
	Stream   *MessageStream
	Producer *MessageStreamProducer
	Consumer *MessageStreamConsumer
}

func (r *Reactor) Run(ctx context.Context) error {
	<-ctx.Done()
	if err := r.Producer.Close(); err != nil {
		return err
	}
	if err := r.Consumer.Close(); err != nil {
		return err
	}

	return nil
}
