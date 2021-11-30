package service

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/kevingentile/chet/pkg/chat"
	"github.com/kevingentile/chet/pkg/user"
)

type RoomService struct {
	bus *cqrs.CommandBus
}

func (rs *RoomService) CreateRoom(host user.UserID) error {
	r := &chat.CreateRoomCmd{
		Host: host,
	}
	return rs.bus.Send(context.TODO(), r)
}

func (rs *RoomService) CreateRoomPublisher(room chat.RoomID) error {
	r := &chat.CreateRoomPublisherCmd{
		ID: room,
	}
	return rs.bus.Send(context.TODO(), r)
}

func (rs *RoomService) DisbandRoom() error {
	return nil
}

func (rs *RoomService) PostMessage() error {
	return nil
}

func NewRoomService(amqpAddress string) (*RoomService, error) {
	logger := watermill.NewStdLogger(false, false)

	commandsPublisher, err := amqp.NewPublisher(amqp.NewDurableQueueConfig(amqpAddress), logger)
	if err != nil {
		return nil, err
	}

	bus, err := cqrs.NewCommandBus(commandsPublisher, func(commandName string) string { return commandName }, cqrs.JSONMarshaler{})
	if err != nil {
		return nil, err
	}
	return &RoomService{
		bus: bus,
	}, nil
}
