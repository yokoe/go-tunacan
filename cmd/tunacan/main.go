package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/yokoe/tunacan/internal/app/tunacan"
)

func main() {
	c := cli.NewCLI("app", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"concat": func() (cli.Command, error) {
			return &tunacan.ConcatCommand{}, nil
		},
		"server": func() (cli.Command, error) {
			return &tunacan.ServerCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
