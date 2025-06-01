package rss

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"feed/feeds"
)

// RSSFeed implements the SocialFeed interface for RSS.
type RSSFeed struct {
	URL string
}

// NewRSSFeed creates a new RSSFeed instance.
func NewRSSFeed(url string) *RSSFeed {
	return &RSSFeed{URL: url}
}

// Fetch retrieves RSS feed items.
func (r *RSSFeed) Fetch() ([]feeds.FeedItem, error) {
	log.Printf("Fetching RSS feed from: %s", r.URL)

	resp, err := http.Get(r.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed from %s: %w", r.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch RSS feed from %s, status code: %d", r.URL, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RSS feed response body from %s: %w", r.URL, err)
	}

	var rssData RSS
	err = xml.Unmarshal(body, &rssData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RSS feed from %s: %w", r.URL, err)
	}

	var items []feeds.FeedItem
	for _, item := range rssData.Channel.Items {
		t, err := time.Parse(time.RFC1123Z, item.PubDate) // Common RSS date format
		if err != nil {
			// Try other common formats if RFC1123Z fails
			t, err = time.Parse(time.RFC3339, item.PubDate)
			if err != nil {
				log.Printf("Warning: Could not parse date '%s' for RSS item from %s: %v", item.PubDate, r.URL, err)
				t = time.Now() // Default to current time if parsing fails
			}
		}

		// Use channel title as username if item author is not available
		username := item.Author
		if username == "" {
			username = rssData.Channel.Title
		}

		// Use link as profile link
		profileLink := item.Link
		if profileLink == "" {
			profileLink = rssData.Channel.Link
		}

		items = append(items, feeds.FeedItem{
			Platform:     "rss",
			PostContent:  item.Title + "\n" + item.Description, // Combine title and description
			Username:     username,
			MediaURL:     nil, // RSS typically doesn't have a direct media_url field like social media
			ProfileLink:  profileLink,
			Timestamp:    t,
			Interactions: 0, // RSS feeds typically don't have interaction counts
		})
	}

	return items, nil
}

// RSS structure for XML unmarshalling
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the RSS channel.
type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Items   []Item   `xml:"item"`
}

// Item represents an individual RSS feed item.
type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Author      string   `xml:"author"`
}
