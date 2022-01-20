package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
)

func init() {
	eh.RegisterAggregate(func(id uuid.UUID) eh.Aggregate {
		return NewUserAggregate(id)
	})
}

var _ = eh.Aggregate(&UserAggregate{})

// TimeNow is a mockable version of time.Now.
var TimeNow = time.Now

const UserAggregateType eh.AggregateType = "User"

type UserAggregate struct {
	// AggregateBase implements most of the eventhorizon.Aggregate interface.
	*events.AggregateBase

	Username  string
	Email     string
	ID        uuid.UUID
	Verified  bool
	CreatedAt int64
}

func NewUserAggregate(id uuid.UUID) *UserAggregate {
	return &UserAggregate{
		AggregateBase: events.NewAggregateBase(UserAggregateType, id),
	}
}

func (a *UserAggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	switch cmd := cmd.(type) {
	case *CreateUser:
		if a.Email != "" || a.Username != "" {
			return errors.New("user already created")
		}
		a.AppendEvent(UserCreatedEvent, &UserCreatedData{
			Username: cmd.Username,
			Email:    cmd.Email,
		}, TimeNow())
	case *VerifyUser:
		if a.Verified {
			return errors.New("user already verified")
		}
		a.AppendEvent(UserVerifiedEvent, nil, TimeNow())
	default:
		return fmt.Errorf("could not handle command: %s", cmd.CommandType())
	}

	return nil
}

func (a *UserAggregate) ApplyEvent(ctx context.Context, event eh.Event) error {
	switch event.EventType() {
	case UserCreatedEvent:
		if data, ok := event.Data().(*UserCreatedData); ok {
			a.Username = data.Username
			a.Email = data.Email
		} else {
			log.Println("invalid UserCreatedData", event.Data())
		}
	case UserVerifiedEvent:
		a.Verified = true
	default:
		return fmt.Errorf("could not apply event: %s", event.EventType())
	}

	return nil
}
