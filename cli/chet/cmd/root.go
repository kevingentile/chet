package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Run() error {
	return rootCmd.Execute()
}
