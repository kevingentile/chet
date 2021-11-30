package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/kevingentile/chet/pkg/service"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type MessageType = string

const (
	MessageCreate MessageType = "CreateRoom"
)

type Message struct {
	Type    MessageType `json:"t"`
	Payload []byte      `json:"p"`
}

type CreateRoomMessage struct {
	ID string `json:"roomId"`
}

func main() {

	messageStream, err := eventstream.NewMessageStream(&eventstream.StreamConfig{
		EnvOptions: stream.NewEnvironmentOptions().
			SetHost("localhost").
			SetPort(5552).
			SetUser("guest").
			SetPassword("guest"),
		StreamOptions: &stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.MB(500),
		},
		StreamName: "chet-messages",
	})
	if err != nil {
		panic(err)
	}

	if err := messageStream.Reset(); err != nil {
		panic(err)
	}

	roomService, _ := service.NewRoomService(messageStream)

	for {
		fmt.Println("creating room")
		if err := roomService.CreateRoom(uuid.New()); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 5)
	}

}
