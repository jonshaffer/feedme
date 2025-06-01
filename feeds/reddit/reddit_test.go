package reddit

import (
	"os"
	"testing"
	"time"

	"feed/feeds"
)

func init() {
	// Dummy usage to ensure feeds import is not removed by linter
	_ = feeds.FeedItem{}
}

func TestRedditFeed_Fetch(t *testing.T) {
	// Mock environment variables for API keys
	t.Setenv("REDDIT_CLIENT_ID", "dummy_client_id")
	t.Setenv("REDDIT_CLIENT_SECRET", "dummy_client_secret")
	t.Setenv("REDDIT_USERNAME", "dummy_username")
	t.Setenv("REDDIT_PASSWORD", "dummy_password")

	redditFeed := NewRedditFeed()
	items, err := redditFeed.Fetch()

	if err != nil {
		t.Fatalf("Fetch returned an error: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Test first item
	expectedPlatform1 := "reddit"
	expectedPostContent1 := "Check out this interesting discussion on r/golang!"
	expectedUsername1 := "GoLover"
	expectedProfileLink1 := "https://reddit.com/user/GoLover"
	expectedInteractions1 := 500

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
	if time.Since(items[0].Timestamp) > 6*time.Hour+5*time.Minute || time.Since(items[0].Timestamp) < 6*time.Hour-5*time.Minute {
		t.Errorf("Item 1 Timestamp is not as expected: %v (expected around %v ago)", items[0].Timestamp, 6*time.Hour)
	}
	if items[0].Interactions != expectedInteractions1 {
		t.Errorf("Item 1 Interactions: Expected %d, got %d", expectedInteractions1, items[0].Interactions)
	}

	// Test second item
	expectedPlatform2 := "reddit"
	expectedPostContent2 := "My thoughts on the latest tech news."
	expectedUsername2 := "NewsReader"
	expectedProfileLink2 := "https://reddit.com/user/NewsReader"
	expectedInteractions2 := 150

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
	if time.Since(items[1].Timestamp) > 28*time.Hour+5*time.Minute || time.Since(items[1].Timestamp) < 28*time.Hour-5*time.Minute {
		t.Errorf("Item 2 Timestamp is not as expected: %v (expected around %v ago)", items[1].Timestamp, 28*time.Hour)
	}
	if items[1].Interactions != expectedInteractions2 {
		t.Errorf("Item 2 Interactions: Expected %d, got %d", expectedInteractions2, items[1].Interactions)
	}
}

func TestRedditFeed_Fetch_MissingAPIKeys(t *testing.T) {
	// Unset environment variables to simulate missing keys
	os.Unsetenv("REDDIT_CLIENT_ID")
	os.Unsetenv("REDDIT_CLIENT_SECRET")
	os.Unsetenv("REDDIT_USERNAME")
	os.Unsetenv("REDDIT_PASSWORD")

	redditFeed := NewRedditFeed()
	_, err := redditFeed.Fetch()

	if err == nil {
		t.Fatalf("Expected an error for missing API keys, got nil")
	}
	expectedError := "Reddit API credentials (client ID, client secret, username, password) not set in environment variables"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
