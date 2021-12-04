package eventstream

import "context"

type EventHandler interface {
	EventName() string
	Handle(ctx context.Context, data []byte) error
}
