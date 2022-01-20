package user

import (
	"context"
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

const UserAggregateType eh.AggregateType = "User"

type Username = string
type Email = string

type UserAggregate struct {
	// AggregateBase implements most of the eventhorizon.Aggregate interface.
	*events.AggregateBase

	Username  string
	Email     string
	ID        uuid.UUID
	Verified  bool
	CreatedAt int64
	// TODO permissions []permission
}

// //TODO
// func validateEmail(email string) (Email, error) {

// }

// func validateUserName(username string) (Username, error) {
//
// }

func NewUserAggregate(id uuid.UUID) *UserAggregate {
	return &UserAggregate{
		AggregateBase: events.NewAggregateBase(UserAggregateType, id),
	}
}

func (a *UserAggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	now := time.Now()
	switch cmd := cmd.(type) {
	case *CreateUser:
		a.AppendEvent(UserCreatedEvent, &UserCreatedData{
			Username:  cmd.Username,
			Email:     cmd.Email,
			CreatedAt: now.Unix(),
		}, now)
		return nil
	case *ConfirmUserCreate:
		a.AppendEvent(UserCreateConfirmedEvent, &ConfirmUserCreate{
			ID:    cmd.ID,
			Email: cmd.Email,
		}, now)
	case *DenyCreateUser:
		a.AppendEvent(UserCreateDeniedEvent, &DenyCreateUser{
			ID: cmd.ID,
		}, now)
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
	// case UserCreateDeniedEvent:
	// 	if data, ok := event.Data().(*UserCreateDeniedData); ok {

	// 	}
	// }
	case UserCreateConfirmedEvent:
		if _, ok := event.Data().(*UserCreateConfirmedData); ok {
			a.CreatedAt = time.Now().Unix()
		} else {
			log.Println("invalid UserCreatedData")
		}

	}
	return nil
}
