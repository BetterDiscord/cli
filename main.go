package main

import (
	"log"

	"github.com/betterdiscord/cli/cmd"
)

func main() {
	log.SetFlags(0)
	cmd.Execute()
}
