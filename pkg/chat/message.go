package chat

import (
	"time"

	"github.com/google/uuid"
)

type MessageSent struct {
	Message     string    `json:"m"`
	Time        time.Time `json:"t"`
	Source      uuid.UUID `json:"s"`
	Destination uuid.UUID `json:"d"`
}

type MessageReceived struct {
	Time        time.Time `json:"t"`
	Source      uuid.UUID `json:"s"`
	Destination uuid.UUID `json:"d"`
}
