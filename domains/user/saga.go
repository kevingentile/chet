package user

import (
	"context"
	"fmt"
	"sync"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/eventhandler/saga"
)

const ResponseSagaType saga.Type = "UserCreated"

type ResponseSaga struct {
	emails   map[string]string
	emailsMu sync.RWMutex
}

func NewResponseSaga() *ResponseSaga {
	return &ResponseSaga{
		emails: make(map[string]string),
	}
}

func (s *ResponseSaga) SagaType() saga.Type {
	return ResponseSagaType
}

func (s *ResponseSaga) RunSaga(ctx context.Context, event eh.Event, h eh.CommandHandler) error {
	switch event.EventType() {
	case Created:
		d, ok := event.Data().(*UserCreatedData)
		if !ok {
			return fmt.Errorf(("saga failed to parse UserCreatedData"))
		}

		s.emailsMu.RLock()
		// user exists already
		_, ok = s.emails[d.Email]
		if ok {
			return h.HandleCommand(ctx, &DenyUserCreate{ID: event.AggregateID()})
		}
		s.emailsMu.RUnlock()

		s.emailsMu.Lock()
		s.emails[d.Email] = d.Username
		s.emailsMu.Unlock()

		return h.HandleCommand(ctx, &ConfirmUserCreate{
			ID: event.AggregateID(),
		})

	}
	return nil
}
