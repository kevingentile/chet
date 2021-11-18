package main

import (
	"github.com/kevingentile/chet/cli/chet/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
