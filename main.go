package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/service"
)

const amqpAddress = "amqp://guest:guest@localhost:5672/"

func main() {
	roomService, err := service.NewRoomService(amqpAddress)
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
