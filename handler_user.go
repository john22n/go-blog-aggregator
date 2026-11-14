package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("expected a second arg: name")
	}
	user := cmd.args[0]
	err := s.Config.SetUser(user)
	if err != nil {
		return err
	}
	fmt.Printf("user set: %s\n", user)
	return nil
}
