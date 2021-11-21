package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(roomCmd)

	roomCmd.AddCommand(newRoomCmd)
}

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Run() error {
	return rootCmd.Execute()
}
