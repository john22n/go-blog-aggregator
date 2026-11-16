package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/john22n/gator/internal/database"
	"github.com/john22n/gator/internal/rssfeed"
)

func handlerAgg(s *state, cmd command) error {
	resp, err := rssfeed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Printf("Feed: %+v\n", resp)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("missing name and url args: %v\n", cmd.name)
	}

	name := cmd.args[0]
	url := cmd.args[1]

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("error getting current user: %v\n", err)
	}

	createdFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v\n", err)
	}

	fmt.Println(createdFeed)

	return nil
}

func handlerGetAllFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting feeds: %v\n", err)
	}

	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.CreatorName)
	}

	return nil
}
