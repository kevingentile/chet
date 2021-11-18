package service

import (
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/kevingentile/chet/pkg/chat"
)

type RoomService struct {
	publisher message.Publisher
}

var (
	publishTopic = "room-events"
)

func (rs *RoomService) CreateRoom() (*chat.Room, error) {
	r := chat.NewRoom()
	roomCreatedEvent := chat.RoomCreatedEvent{
		Room: *r,
		Time: time.Now(),
	}
	payload, err := json.Marshal(roomCreatedEvent)
	if err != nil {
		return nil, err
	}

	if err := rs.publish(payload); err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *RoomService) publish(payload message.Payload) error {
	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(watermill.NewShortUUID(), msg)
	if err := rs.publisher.Publish(publishTopic, msg); err != nil {
		return err
	}
	return nil
}

func (rs *RoomService) DisbandRoom() {

}

func NewRoomService(brokers []string) *RoomService {
	p := createPublisher(brokers)
	return &RoomService{
		publisher: p,
	}
}

func createPublisher(brokers []string) message.Publisher {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(
			true,  // TODO flag debug
			false, // trace
		),
	)
	if err != nil {
		panic(err)
	}

	return kafkaPublisher
}
