package service

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/kevingentile/chet/pkg/chat"
	"github.com/kevingentile/chet/pkg/user"
	"google.golang.org/protobuf/proto"
)

type RoomService struct {
	publisher message.Publisher
}

var (
	publishTopic = "room-events"
)

func (rs *RoomService) CreateRoom(host user.UserID) error {
	r := &chat.CreateRoom{
		HostID: host.String(),
	}

	payload, err := proto.Marshal(r)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(watermill.NewShortUUID(), msg)
	if err := rs.publisher.Publish(publishTopic, msg); err != nil {
		return err
	}
	return nil
}

func (rs *RoomService) DisbandRoom() error {
	return nil
}

func (rs *RoomService) PostMessage() error {
	return nil
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
