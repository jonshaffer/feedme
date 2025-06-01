package credly

import (
	"testing"
	"time"
)

func TestNewCredlyFeed(t *testing.T) {
	feed := NewCredlyFeed()
	if feed == nil {
		t.Error("NewCredlyFeed returned nil")
	}
}

func TestCredlyFeed_Fetch(t *testing.T) {
	feed := NewCredlyFeed()
	items, err := feed.Fetch()
	if err != nil {
		t.Fatalf("Fetch returned an error: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	if len(items) > 0 {
		item := items[0]
		if item.Platform != "credly" {
			t.Errorf("Expected platform 'credly', got '%s'", item.Platform)
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
