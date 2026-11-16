package rssfeed

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	//setup reusable client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// build request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	// always close resp body to prevent source leaks
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(respBody, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}

	sanitizeFeed(&feed)

	return &feed, nil
}

func sanitizeFeed(feed *RSSFeed) {
	channel := feed.Channel

	channel.Title = html.UnescapeString(channel.Title)
	channel.Description = html.UnescapeString(channel.Description)

	for i := range channel.Item {
		item := &channel.Item[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
}
