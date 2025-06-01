package strava

import (
	"log"
	"time"

	"feed/feeds"
)

// StravaFeed implements the SocialFeed interface for Strava.
type StravaFeed struct {
	// Add any Strava-specific configuration here, e.g., API client
}

// NewStravaFeed creates a new StravaFeed instance.
func NewStravaFeed() *StravaFeed {
	return &StravaFeed{}
}

// Fetch retrieves feed items from Strava.
func (sf *StravaFeed) Fetch() ([]feeds.FeedItem, error) {
	log.Println("Simulating Strava feed fetch...")
	// In a real implementation, this would involve calling the Strava API.
	// For now, return a dummy item or an empty slice.
	return []feeds.FeedItem{
		{
			Platform:     "strava",
			PostContent:  "Just completed a 10k run!",
			Username:     "strava_user",
			MediaURL:     nil,
			ProfileLink:  "https://www.strava.com/athletes/strava_user",
			Timestamp:    time.Now(),
			Interactions: 15,
		},
	}, nil
}
