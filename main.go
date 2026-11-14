package main

import (
	"fmt"
	"log"

	"os"

	"github.com/john22n/gator/internal/config"
)

type state struct {
	Config *config.Config
}

func main() {
	Args := os.Args
	if len(Args) < 2 {
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	s := state{Config: &cfg}

	cmd := command{name: Args[1], args: Args[2:]}
	cmds := commands{registeredCommands: map[string]func(*state, command) error{}}
	cmds.register("login", handlerLogin)

	if err := cmds.run(&s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
