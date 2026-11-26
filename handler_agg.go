package main

import (
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("missing time arg")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error parsing time args: %w", err)
	}

	ticker := time.NewTicker(timeBetweenReqs)

	fmt.Printf("Colleting feed every: %v", timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
