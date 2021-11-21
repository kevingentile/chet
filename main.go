package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/service"
)

var amqpAddress = "amqp://guest:guest@localhost:5672/"

func main() {
	roomService, err := service.NewRoomService(amqpAddress)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("Publishing create room...")
		roomService.CreateRoom(uuid.New())
		time.Sleep(1 * time.Second)
	}
}
