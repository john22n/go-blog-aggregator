package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/john22n/gator/internal/database"
	"time"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected a second arg: %s", cmd.name)
	}

	userArg := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), userArg)
	if err != nil {
		return fmt.Errorf("coudl not find user: %s", err)
	}
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldnt set user: %w", err)
	}
	fmt.Printf("successful login: %s\n", userArg)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected a second arg: %s", cmd.name)
	}
	username := cmd.args[0]
	createdUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: username})
	if err != nil {
		return fmt.Errorf("couldnt create user: %s", err)
	}
	err = s.Config.SetUser(createdUser.Name)
	if err != nil {
		return fmt.Errorf("couldnt set the created user: %s", err)
	}
	printUser(createdUser)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID: %v\n", user.ID)
	fmt.Printf(" * Name: %v\n", user.Name)
}
