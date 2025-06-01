package goodreads

import (
	"testing"
	"time"
)

func TestNewGoodreadsFeed(t *testing.T) {
	feed := NewGoodreadsFeed()
	if feed == nil {
		t.Error("NewGoodreadsFeed returned nil")
	}
}

func TestGoodreadsFeed_Fetch(t *testing.T) {
	feed := NewGoodreadsFeed()
	items, err := feed.Fetch()
	if err != nil {
		t.Fatalf("Fetch returned an error: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	if len(items) > 0 {
		item := items[0]
		if item.Platform != "goodreads" {
			t.Errorf("Expected platform 'goodreads', got '%s'", item.Platform)
		}
		if item.PostContent == "" {
			t.Error("PostContent is empty")
		}
		if item.Username == "" {
			t.Error("Username is empty")
		}
		if item.ProfileLink == "" {
			t.Error("ProfileLink is empty")
		}
		if item.Timestamp.IsZero() {
			t.Error("Timestamp is zero")
		}
		// Check if timestamp is recent (e.g., within the last minute)
		if time.Since(item.Timestamp) > time.Minute {
			t.Errorf("Timestamp is not recent: %v", item.Timestamp)
		}
	}
}
