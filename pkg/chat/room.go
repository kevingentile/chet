package chat

import (
	"time"

	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/user"
)

type RoomID = uuid.UUID

const (
	CreateRoomCmd          = "CreateRoomCmd"
	CreateRoomPublisherCmd = "CreateRoomPublisherCmd"

	RoomCreatedEvent = "RoomCreatedEvent"
)

type Room struct {
	Host  user.UserID   `json:"host"`
	Peers []user.UserID `json:"peers"`
	ID    RoomID        `json:"roomId"`
}

type CreateRoom struct {
	Host user.UserID `json:"host"`
}

type RoomCreated struct {
	Room Room      `json:"room"`
	Time time.Time `json:"createdAt"`
}

type RoomDisbanded struct {
	ID   RoomID    `json:"rid"`
	Time time.Time `json:"td"`
}

type RoomUserAdded struct {
	ID     RoomID      `json:"rid"`
	UserID user.UserID `json:"uid"`
}

type CreateRoomPublisher struct {
	ID RoomID `json:"rid"`
}

type RoomPublisherCreated struct {
	ID        RoomID `json:"rid"`
	TopicName string `json:"topicName"`
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
	ID        string    `bson:"id"`
	Host      string    `bson:"host"`
	Peers     []string  `bson:"peers"`
	CreatedAt time.Time `bson:"createdAt"`
}

func NewRoomView(room Room) *RoomView {
	peers := make([]string, len(room.Peers))
	for i, p := range room.Peers {
		peers[i] = p.String()
	}
	return &RoomView{
		ID:        room.ID.String(),
		Host:      room.Host.String(),
		Peers:     peers,
		CreatedAt: time.Now().UTC(),
	}
}
