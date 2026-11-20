package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/john22n/gator/internal/database"
)

func handlerAddFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("missing url argument")
	}

	urlArg := cmd.args[0]

	currentFeed, err := s.db.GetFeedByUrl(context.Background(), urlArg)
	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    currentFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("issue creating follow: %w", err)
	}

	printFollow(follow)

	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
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

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("missing feed arg")
	}

	feedURL := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	unfollowArgs := database.DeleteFeedFollowFromUserParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}

	err = s.db.DeleteFeedFollowFromUser(context.Background(), unfollowArgs)
	if err != nil {
		return fmt.Errorf("error deleteing feed from user: %w", err)
	}
	fmt.Println("successfully removed feed from user")

	return nil
}
