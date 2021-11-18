package chat

import (
	"time"

	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/user"
)

type RoomID = uuid.UUID

type Room struct {
	Peers []user.UserID `json:"p"`
	ID    RoomID        `json:"rid"`
}

type RoomCreatedEvent struct {
	Room Room      `json:"r"`
	Time time.Time `json:"tc"`
}

type RoomDisbandedEvent struct {
	ID   RoomID    `json:"rid"`
	Time time.Time `json:"td"`
}

type RoomUserAddedEvent struct {
	ID     RoomID      `json:"rid"`
	UserID user.UserID `json:"uid"`
}

func NewRoom() *Room {
	return &Room{
		Peers: []user.UserID{},
		ID:    uuid.New(),
	}
}

func (r *Room) AddUser(userID user.UserID) {
	r.Peers = append(r.Peers, userID)
}

func (r *Room) RemoveUser(userID user.UserID) {
	for index, uid := range r.Peers {
		if uid == userID {
			r.Peers = append(r.Peers[:index], r.Peers[index+1])
			break
		}
	}
}

type RoomServicer interface {
	CreateRoom() (*Room, error)
	DisbandRoom() error
	PostMessage() error
}
