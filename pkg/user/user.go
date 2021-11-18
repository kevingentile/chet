package user

import (
	"time"

	"github.com/google/uuid"
)

type UserID = uuid.UUID

type User struct {
	Alias string `json:"alias"`
	ID    UserID `json:"uid"`
	// Name  string    `json:"name"` //TODO associate personal info with keymap
}

type UserCreatedEvent struct {
	User User      `json:"user"` //TODO likely don't wan't the whole user
	Time time.Time `json:"t"`
}

func NewUser(alias string) *User {
	return &User{
		Alias: alias,
		ID:    uuid.New(),
	}
}

type UserServicer interface {
	CreateUser() *User
}
