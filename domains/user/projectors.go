package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	"github.com/looplab/eventhorizon/uuid"
)

type User struct {
	ID        uuid.UUID `bson:"_id"`
	Version   int
	Username  string
	Email     string
	Verified  bool
	CreatedAt int64
}

var _ = eh.Entity(&User{})
var _ = eh.Versionable(&User{})

func (u *User) EntityID() uuid.UUID {
	return u.ID
}

func (u *User) AggregateVersion() int {
	return u.Version
}

type UserProjector struct{}

func NewUserProjector() *UserProjector {
	return &UserProjector{}
}

func (p *UserProjector) ProjectorType() projector.Type {
	return projector.Type(UserAggregateType.String())
}

func (p *UserProjector) Project(ctx context.Context, event eh.Event, entity eh.Entity) (eh.Entity, error) {
	u, ok := entity.(*User)
	if !ok {
		return nil, errors.New("model is of incorrect type")
	}

	switch event.EventType() {
	case UserCreatedEvent:
		data, ok := event.Data().(*UserCreatedData)
		if !ok {
			return nil, fmt.Errorf("projector: invalid event data type: %v", event.Data())
		}
		u.ID = event.AggregateID()
		u.Username = data.Username
		u.Email = data.Email
		u.CreatedAt = time.Now().UTC().Unix()

	case UserVerifiedEvent:
		u.Verified = true
	default:
		return nil, fmt.Errorf("could not handle event: %s", event)
	}

	u.Version++
	return u, nil
}
