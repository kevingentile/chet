package chat

import (
	"time"

	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/user"
)

type RoomID = uuid.UUID

type Room struct {
	Host  user.UserID   `json:"host"`
	Peers []user.UserID `json:"peers"`
	ID    RoomID        `json:"roomId"`
}

type CreateRoomCmd struct {
	Host user.UserID `json:"host"`
}

type RoomCreatedEvent struct {
	Room Room      `json:"room"`
	Time time.Time `json:"createdAt"`
}

type RoomDisbandedEvent struct {
	ID   RoomID    `json:"rid"`
	Time time.Time `json:"td"`
}

type RoomUserAddedEvent struct {
	ID     RoomID      `json:"rid"`
	UserID user.UserID `json:"uid"`
}

func NewRoom(host user.UserID) *Room {
	return &Room{
		Host:  host,
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

type RoomView struct {
	Host      user.UserID   `bson:"host"`
	Peers     []user.UserID `bson:"peers"`
	ID        RoomID        `bson:"_id"`
	CreatedAt time.Time     `bson:"createdAt"`
}

func NewRoomView(room *Room) *RoomView {
	return &RoomView{
		Host:      room.Host,
		Peers:     room.Peers,
		ID:        room.ID,
		CreatedAt: time.Now().UTC(),
	}
}
