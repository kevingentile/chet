package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kevingentile/chet/pkg/chat"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/kevingentile/chet/pkg/service"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type Message struct {
	// Type string `json:"type"`
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type RoomRelay struct {
	RoomID string
	// RoomName     string
	Subscription *eventstream.MessageStreamConsumer
	Publisher    *eventstream.MessageStreamProducer
}

var roomRelays []*RoomRelay = make([]*RoomRelay, 0)
var messageChannels []<-chan *message.Message = []<-chan *message.Message{}

func main() {
	router := gin.Default()
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return r.Host == "localhost:8000"
	}

	go runMessageBroker(context.TODO())
	// go handleRoomSubscriptions()

	router.GET("/ws", func(c *gin.Context) {
		roomID := c.Query("roomId")
		if roomID == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		var relay *RoomRelay
		for _, r := range roomRelays {
			if r.RoomID == roomID {
				relay = r
				break
			}
		}
		if relay == nil {
			c.Status(http.StatusNotFound)
			return
		}
		con, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer con.Close()

		for {
			m, err := json.Marshal(Message{
				ID:        roomID,
				Timestamp: time.Now().Unix(),
			})
			if err != nil {
				log.Println(err)
			}
			if err := con.WriteMessage(websocket.TextMessage, m); err != nil {
				break
			}
			time.Sleep(2 * time.Second)
		}
	})
	router.Run("localhost:8000")
}

func runMessageBroker(ctx context.Context) {
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

	eventStream, err := eventstream.NewEventStream(messageStream)
	if err != nil {
		panic(err)
	}

	roomService, err := service.NewRoomService(messageStream)
	if err != nil {
		panic(err)
	}

	roomCreatedHandler := &roomCreatedHandler{
		roomService: roomService,
	}

	eventStream.AddEventHandler(roomCreatedHandler)

	if err := eventStream.Run(ctx); err != nil {
		panic(err)
	}
}

type roomCreatedHandler struct {
	roomService *service.RoomService
}

func (h *roomCreatedHandler) EventName() string {
	return chat.RoomCreatedEvent
}

func (h *roomCreatedHandler) Handle(ctx context.Context, data []byte) error {
	e := &chat.RoomCreated{}
	if err := json.Unmarshal(data, e); err != nil {
		return err
	}

	fmt.Println("Handling room created event")

	return nil
}
