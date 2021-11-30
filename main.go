package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
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

	reader := bufio.NewReader(os.Stdin)

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

	messageProducer, err := messageStream.NewProducer(&eventstream.MessageStreamProducerConfig{})
	if err != nil {
		panic(err)
	}
	go func() {
		for i := 0; i < 2000; i++ {
			cr := &CreateRoomMessage{ID: uuid.NewString()}
			crp, err := json.Marshal(cr)
			if err != nil {
				panic(err)
			}
			msg := &Message{
				Type:    MessageCreate,
				Payload: crp,
			}

			payload, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			err = messageProducer.Send(amqp.NewMessage(payload))
			CheckErr(err)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	handleMessages := func(consumerContext stream.ConsumerContext, msg *amqp.Message) {
		for _, b := range msg.Data {
			cmsg := &Message{}
			if err := json.Unmarshal(b, cmsg); err != nil {
				panic(err)
			}
			switch cmsg.Type {
			case MessageCreate:
				crm := &CreateRoomMessage{}
				if err := json.Unmarshal(cmsg.Payload, crm); err != nil {
					panic(err)
				}
				fmt.Println(crm)
			default:
				// panic("unmached message")
			}
			// fmt.Println(string(b))

		}
	}

	messageConsumer, err := messageStream.NewConsumer(&eventstream.MessageStreamConsumerConfig{
		Options: stream.NewConsumerOptions().
			SetConsumerName("my_consumer").SetOffset(stream.OffsetSpecification{}.First()),
		MessageHandler: handleMessages,
		Offset:         stream.OffsetSpecification{}.First(),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Press any key to stop ")
	_, _ = reader.ReadString('\n')
	err = messageProducer.Close()
	CheckErr(err)
	err = messageConsumer.Close()
	CheckErr(err)
	err = messageStream.Delete()
	CheckErr(err)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("%s ", err)
		os.Exit(1)
	}
}
