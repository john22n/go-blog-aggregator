package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/john22n/gator/internal/database"
)

func handlerAddFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("missing url argument")
	}

	urlArg := cmd.args[0]

	currentFeed, err := s.db.GetFeedByUrl(context.Background(), urlArg)
	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("error getting user details from db: %w", err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    currentFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("issue creating follow: %w", err)
	}

	printFollow(follow)

	return nil
}

func handlerListFeedFollows(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("issue getting user info: %v", err)
	}
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting follows: %w\n", err)
	}

	for _, follow := range follows {
		fmt.Printf("feed name: %v\n", follow.FeedName)
	}

	return nil
}

func printFollow(follow database.CreateFeedFollowRow) {
	fmt.Printf(" * Name: %v\n", follow.UserName)
	fmt.Printf(" * Feed Name: %v\n", follow.FeedName)
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User: %s", username)
	fmt.Printf("* Feed: %s", feedname)
}
