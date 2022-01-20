package user

import (
	eh "github.com/looplab/eventhorizon"
)

const (
	UserCreatedEvent  eh.EventType = "UserCreated"
	UserVerifiedEvent eh.EventType = "UserVerified"
)

func init() {
	eh.RegisterEventData(UserCreatedEvent, func() eh.EventData {
		return &UserCreatedData{}
	})
}

type UserCreatedData struct {
	Username string `bson:"username"`
	Email    string `bson:"email"`
}
