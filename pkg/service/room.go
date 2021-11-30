package service

import (
	"encoding/json"

	"github.com/kevingentile/chet/pkg/chat"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/kevingentile/chet/pkg/user"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
)

type RoomService struct {
	producer *eventstream.MessageStreamProducer
}

func (rs *RoomService) CreateRoom(host user.UserID) error {
	r := &chat.CreateRoom{
		Host: host,
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return err
	}

	msg := amqp.NewMessage(payload)
	msg.Properties = &amqp.MessageProperties{Subject: chat.CreateRoomCmd}
	return rs.producer.Send(msg)
}

func (rs *RoomService) CreateRoomPublisher(room chat.RoomID) error {
	r := &chat.CreateRoomPublisher{
		ID: room,
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return err
	}
	msg := amqp.NewMessage(payload)
	msg.Properties = &amqp.MessageProperties{Subject: chat.CreateRoomPublisherCmd}
	return rs.producer.Send(msg)
}

func (rs *RoomService) DisbandRoom() error {
	return nil
}

func (rs *RoomService) PostMessage() error {
	return nil
}

func NewRoomService(stream *eventstream.MessageStream) (*RoomService, error) {
	messageProducer, err := stream.NewProducer(&eventstream.MessageStreamProducerConfig{})
	if err != nil {
		panic(err)
	}

	return &RoomService{
		producer: messageProducer,
	}, nil
}
