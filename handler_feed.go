package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/john22n/gator/internal/database"
	"github.com/john22n/gator/internal/rssfeed"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("missing name and url args: %v\n", cmd.name)
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v\n", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error creating feed: %v\n", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerGetAllFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting feeds: %v\n", err)
	}

	if len(feeds) == 0 {
		fmt.Println("couldnt get feeds: %w\n", err)
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldnt get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("==================================================")
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:               %s\n", feed.ID)
	fmt.Printf("* Created: 					%v\n", feed.CreatedAt)
	fmt.Printf("* Updated: 					%v\n", feed.UpdatedAt)
	fmt.Printf("* Name: 						%s\n", feed.Name)
	fmt.Printf("* URL:							%s\n", feed.Url)
	fmt.Printf("* User:							%s\n", user.Name)
}

func scrapeFeeds(s *state) {
	//get next feed
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("error getting next feed: %w", err)
	}
	log.Println("found a feed to fetch!")
	scrapeFeed(s.db, nextFeed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldnt makr feed %s fetched: %v", feed.Name, err)
		return
	}
	feedData, err := rssfeed.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("couldnt collect feed %s: %v", feed.Name, err)
	}

	for _, feedItem := range feedData.Channel.Item {
		parsedPubDate, err := time.Parse(time.RFC1123Z, feedItem.PubDate)
		if err != nil {
			log.Printf("error parsing pubDate: %s", err)
		}

		savedPost, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     feedData.Channel.Title,
			Description: sql.NullString{
				String: feedData.Channel.Description,
				Valid:  feedData.Channel.Description != "",
			},
			PublishedAt: sql.NullTime{
				Time:  parsedPubDate,
				Valid: true,
			},
			FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				continue
			} else {
				log.Printf("error creating post from feed: %s", err)
			}
		}
		log.Printf("Post created from feed: %v", savedPost)
	}

}
