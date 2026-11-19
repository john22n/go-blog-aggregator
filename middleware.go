package main

import (
	//	"context"
	//	"fmt"

	"github.com/john22n/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return nil
}
