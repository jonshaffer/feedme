package feeds

import "time"

// FeedItem represents a standardized social media post or RSS item.
type FeedItem struct {
	Platform     string    `json:"platform"`
	PostContent  string    `json:"post_content"`
	Username     string    `json:"username"`
	MediaURL     *string   `json:"media_url"` // Use pointer for nullable string
	ProfileLink  string    `json:"profile_link"`
	Timestamp    time.Time `json:"timestamp"`
	Interactions int       `json:"interactions"`
	Permalink    string    `json:"permalink,omitempty"` // URL to the individual item's JSON file, if generated
}

// SocialFeed defines the interface for fetching social media feed items.
type SocialFeed interface {
	Fetch() ([]FeedItem, error)
}
