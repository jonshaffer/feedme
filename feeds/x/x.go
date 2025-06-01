package x

import (
	"fmt"
	"log"
	"os"
	"time"

	"feed/feeds"
)

// XFeed implements the SocialFeed interface for X (formerly Twitter).
type XFeed struct{}

// NewXFeed creates a new XFeed instance.
func NewXFeed() *XFeed {
	return &XFeed{}
}

// Fetch retrieves X feed items.
func (x *XFeed) Fetch() ([]feeds.FeedItem, error) {
	// In a real scenario, this would involve X API calls.
	// For now, we'll simulate fetching data.

	apiKey := os.Getenv("X_API_KEY")
	apiSecret := os.Getenv("X_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("X API keys not set in environment variables")
	}

	log.Println("Simulating X feed fetch...")

	// Simulate fetching a few items
	items := []feeds.FeedItem{
		{
			Platform:     "x",
			PostContent:  "Just shared a new article on Go concurrency! #golang #programming",
			Username:     "GoDev",
			MediaURL:     nil,
			ProfileLink:  "https://x.com/godev",
			Timestamp:    time.Now().Add(-3 * time.Hour),
			Interactions: 200,
		},
		{
			Platform:     "x",
			PostContent:  "Excited for the upcoming tech conference!",
			Username:     "TechEnthusiast",
			MediaURL:     nil,
			ProfileLink:  "https://x.com/techenthusiast",
			Timestamp:    time.Now().Add(-18 * time.Hour),
			Interactions: 90,
		},
	}

	return items, nil
}
