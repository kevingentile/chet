package cmd

import (
	"github.com/google/uuid"
	"github.com/kevingentile/chet/pkg/eventstream"
	"github.com/kevingentile/chet/pkg/service"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"github.com/spf13/cobra"
)

var roomCmd = &cobra.Command{
	Use: "room",
	Run: func(cmd *cobra.Command, args []string) { cmd.Help() },
}

var newRoomCmd = &cobra.Command{
	Use:   "new",
	Short: "generate a random room",
	Run: func(cmd *cobra.Command, args []string) {
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
		roomService, err := service.NewRoomService(messageStream)
		if err != nil {
			panic(err)
		}
		roomService.CreateRoom(uuid.New())
	},
}
