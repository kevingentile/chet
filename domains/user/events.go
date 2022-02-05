package user

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

const (
	CreatedEvent         eh.EventType = "user:created"
	CreateConfirmedEvent eh.EventType = "user:createConfirmed"
	VerifiedEvent        eh.EventType = "user:verified"
)

func init() {
	eh.RegisterEventData(CreatedEvent, func() eh.EventData {
		return &UserCreatedData{}
	})

	eh.RegisterEventData(CreateConfirmedEvent, func() eh.EventData {
		return &CreateConfirmedData{}
	})

	eh.RegisterEventData(VerifiedEvent, func() eh.EventData {
		return &UserVerifiedData{}
	})

}

type UserCreatedData struct {
	Username string `bson:"username"`
	Email    string `bson:"email"`
}

type CreateConfirmedData struct {
	ID uuid.UUID `bson:"id"`
}

type UserVerifiedData struct {
	ID    uuid.UUID `bson:"id"`
	Email string    `bson:"email"`
}
