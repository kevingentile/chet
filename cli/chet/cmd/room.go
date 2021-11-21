package cmd

import (
	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/service"
	"github.com/spf13/cobra"
)

//TODO move to config
var amqpAddress = "amqp://guest:guest@localhost:5672/"

var roomCmd = &cobra.Command{
	Use: "room",
	Run: func(cmd *cobra.Command, args []string) { cmd.Help() },
}

var newRoomCmd = &cobra.Command{
	Use:   "new",
	Short: "generate a random room",
	Run: func(cmd *cobra.Command, args []string) {
		roomService, err := service.NewRoomService(amqpAddress)
		if err != nil {
			panic(err)
		}
		roomService.CreateRoom(uuid.New())
	},
}
