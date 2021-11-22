package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/infrastructure"
	"github.com/kevingentile/chet/pkg/service"
)

const amqpAddress = "amqp://guest:guest@localhost:5672/"
const mongoUri = "mongodb://root:example@localhost:27017/?maxPoolSize=20&w=majority"

func main() {
	eventStore, err := infrastructure.NewDefaultMongoStore(mongoUri)
	if err != nil {
		panic(err)
	}
	roomService, err := service.NewRoomService(amqpAddress, eventStore)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("Publishing create room...")
		if err := roomService.CreateRoom(uuid.New()); err != nil {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
}
