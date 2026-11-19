package main

import (
	"context"
	"fmt"
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
