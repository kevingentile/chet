package user

import (
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/uuid"
)

const (
	UserCreatedEvent         eh.EventType = "UserCreated"
	UserVerifiedEvent        eh.EventType = "UserVerified"
	UserCreateDeniedEvent    eh.EventType = "UserCreateDenied"
	UserCreateConfirmedEvent eh.EventType = "UserCreateConfirmed"
)

func init() {
	eh.RegisterEventData(UserCreatedEvent, func() eh.EventData {
		return &UserCreatedData{}
	})
	eh.RegisterEventData(UserCreateDeniedEvent, func() eh.EventData {
		return &UserCreateDeniedData{}
	})
	eh.RegisterEventData(UserCreateConfirmedEvent, func() eh.EventData {
		return &UserCreateConfirmedData{}
	})

}

type UserCreatedData struct {
	Username  string `bson:"username"`
	Email     string `bson:"email"`
	CreatedAt int64  `bson:"createdAt"`
}

type UserCreateDeniedData struct {
	ID uuid.UUID `bson:"id"`
}
type UserCreateConfirmedData struct {
	ID    uuid.UUID `bson:"id"`
	Email string    `bson:"email"`
}
