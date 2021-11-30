package eventstream

import "context"

type Projector struct {
	Stream   *MessageStream
	Consumer *MessageStreamConsumer
}

func (r *Projector) Run(ctx context.Context) error {
	<-ctx.Done()
	if err := r.Consumer.Close(); err != nil {
		return err
	}

	return nil
}
