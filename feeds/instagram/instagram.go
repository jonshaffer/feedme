package instagram

import (
	"fmt"
	"log"
	"os"
	"time"

	"feed/feeds"
)

// InstagramFeed implements the SocialFeed interface for Instagram.
type InstagramFeed struct{}

// NewInstagramFeed creates a new InstagramFeed instance.
func NewInstagramFeed() *InstagramFeed {
	return &InstagramFeed{}
}

// Fetch retrieves Instagram feed items.
func (i *InstagramFeed) Fetch() ([]feeds.FeedItem, error) {
	// In a real scenario, this would involve Instagram API calls.
	// For now, we'll simulate fetching data.

	apiKey := os.Getenv("INSTAGRAM_API_KEY")
	apiSecret := os.Getenv("INSTAGRAM_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("Instagram API keys not set in environment variables")
	}

	log.Println("Simulating Instagram feed fetch...")

	// Simulate fetching a few items
	items := []feeds.FeedItem{
		{
			Platform:     "instagram",
			PostContent:  "Beautiful sunset views from the beach! #travel #photography",
			Username:     "TravelBug",
			MediaURL:     stringPtr("https://instagram.com/p/sunset.jpg"),
			ProfileLink:  "https://instagram.com/travelbug",
			Timestamp:    time.Now().Add(-5 * time.Hour),
			Interactions: 350,
		},
		{
			Platform:     "instagram",
			PostContent:  "New recipe alert! Delicious homemade pasta.",
			Username:     "FoodieChef",
			MediaURL:     stringPtr("https://instagram.com/p/pasta.jpg"),
			ProfileLink:  "https://instagram.com/foodiechef",
			Timestamp:    time.Now().Add(-20 * time.Hour),
			Interactions: 280,
		},
	}

	return items, nil
}

// Helper function to return a pointer to a string
func stringPtr(s string) *string {
	return &s
}
