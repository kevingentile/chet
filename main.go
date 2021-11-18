package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/service"
)

func main() {
	brokers := []string{"127.0.0.1:9092"}
	rs := service.NewRoomService(brokers)
	// user := user.NewUser("test")
	err := rs.CreateRoom(uuid.New())
	if err != nil {
		panic(err)
	}
	fmt.Println("created room...")

}
