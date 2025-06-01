package linkedin

import (
	"os"
	"testing"
	"time"

	"feed/feeds" // Required for feeds.FeedItem type
)

func init() {
	// Dummy usage to ensure feeds import is not removed by linter
	_ = feeds.FeedItem{}
}

func TestLinkedInFeed_Fetch(t *testing.T) {
	// Mock environment variables for API keys
	t.Setenv("LINKEDIN_API_KEY", "dummy_key")
	t.Setenv("LINKEDIN_API_SECRET", "dummy_secret")

	liFeed := NewLinkedInFeed()
	items, err := liFeed.Fetch()

	if err != nil {
		t.Fatalf("Fetch returned an error: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Test first item
	expectedPlatform1 := "linkedin"
	expectedPostContent1 := "Excited to share my latest project! #golang #softwaredevelopment"
	expectedUsername1 := "Jane Doe"
	expectedProfileLink1 := "https://linkedin.com/in/janedoe"
	expectedInteractions1 := 120

	if items[0].Platform != expectedPlatform1 {
		t.Errorf("Item 1 Platform: Expected %s, got %s", expectedPlatform1, items[0].Platform)
	}
	if items[0].PostContent != expectedPostContent1 {
		t.Errorf("Item 1 PostContent: Expected %s, got %s", expectedPostContent1, items[0].PostContent)
	}
	if items[0].Username != expectedUsername1 {
		t.Errorf("Item 1 Username: Expected %s, got %s", expectedUsername1, items[0].Username)
	}
	if items[0].MediaURL != nil {
		t.Errorf("Item 1 MediaURL: Expected nil, got %v", *items[0].MediaURL)
	}
	if items[0].ProfileLink != expectedProfileLink1 {
		t.Errorf("Item 1 ProfileLink: Expected %s, got %s", expectedProfileLink1, items[0].ProfileLink)
	}
	// Check timestamp is recent (within a reasonable margin)
	// The simulated timestamp is time.Now().Add(-2 * time.Hour)
	// We check if it's within a small window around that expected time.
	if time.Since(items[0].Timestamp) > 2*time.Hour+5*time.Minute || time.Since(items[0].Timestamp) < 2*time.Hour-5*time.Minute {
		t.Errorf("Item 1 Timestamp is not as expected: %v (expected around %v ago)", items[0].Timestamp, 2*time.Hour)
	}
	if items[0].Interactions != expectedInteractions1 {
		t.Errorf("Item 1 Interactions: Expected %d, got %d", expectedInteractions1, items[0].Interactions)
	}

	// Test second item
	expectedPlatform2 := "linkedin"
	expectedPostContent2 := "Great discussion on microservices architecture today."
	expectedUsername2 := "John Smith"
	expectedProfileLink2 := "https://linkedin.com/in/johnsmith"
	expectedInteractions2 := 85

	if items[1].Platform != expectedPlatform2 {
		t.Errorf("Item 2 Platform: Expected %s, got %s", expectedPlatform2, items[1].Platform)
	}
	if items[1].PostContent != expectedPostContent2 {
		t.Errorf("Item 2 PostContent: Expected %s, got %s", expectedPostContent2, items[1].PostContent)
	}
	if items[1].Username != expectedUsername2 {
		t.Errorf("Item 2 Username: Expected %s, got %s", expectedUsername2, items[1].Username)
	}
	if items[1].MediaURL != nil {
		t.Errorf("Item 2 MediaURL: Expected nil, got %v", *items[1].MediaURL)
	}
	if items[1].ProfileLink != expectedProfileLink2 {
		t.Errorf("Item 2 ProfileLink: Expected %s, got %s", expectedProfileLink2, items[1].ProfileLink)
	}
	// Check timestamp is recent (within a reasonable margin)
	// The simulated timestamp is time.Now().Add(-24 * time.Hour)
	if time.Since(items[1].Timestamp) > 24*time.Hour+5*time.Minute || time.Since(items[1].Timestamp) < 24*time.Hour-5*time.Minute {
		t.Errorf("Item 2 Timestamp is not as expected: %v (expected around %v ago)", items[1].Timestamp, 24*time.Hour)
	}
	if items[1].Interactions != expectedInteractions2 {
		t.Errorf("Item 2 Interactions: Expected %d, got %d", expectedInteractions2, items[1].Interactions)
	}
}

func TestLinkedInFeed_Fetch_MissingAPIKeys(t *testing.T) {
	// Unset environment variables to simulate missing keys
	os.Unsetenv("LINKEDIN_API_KEY")
	os.Unsetenv("LINKEDIN_API_SECRET")

	liFeed := NewLinkedInFeed()
	_, err := liFeed.Fetch()

	if err == nil {
		t.Fatalf("Expected an error for missing API keys, got nil")
	}
	expectedError := "LinkedIn API keys not set in environment variables"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
