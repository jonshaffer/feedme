package instagram

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

func TestInstagramFeed_Fetch(t *testing.T) {
	// Mock environment variables for API keys
	t.Setenv("INSTAGRAM_API_KEY", "dummy_key")
	t.Setenv("INSTAGRAM_API_SECRET", "dummy_secret")

	instagramFeed := NewInstagramFeed()
	items, err := instagramFeed.Fetch()

	if err != nil {
		t.Fatalf("Fetch returned an error: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Test first item
	expectedPlatform1 := "instagram"
	expectedPostContent1 := "Beautiful sunset views from the beach! #travel #photography"
	expectedUsername1 := "TravelBug"
	expectedMediaURL1 := "https://instagram.com/p/sunset.jpg"
	expectedProfileLink1 := "https://instagram.com/travelbug"
	expectedInteractions1 := 350

	if items[0].Platform != expectedPlatform1 {
		t.Errorf("Item 1 Platform: Expected %s, got %s", expectedPlatform1, items[0].Platform)
	}
	if items[0].PostContent != expectedPostContent1 {
		t.Errorf("Item 1 PostContent: Expected %s, got %s", expectedPostContent1, items[0].PostContent)
	}
	if items[0].Username != expectedUsername1 {
		t.Errorf("Item 1 Username: Expected %s, got %s", expectedUsername1, items[0].Username)
	}
	if items[0].MediaURL == nil || *items[0].MediaURL != expectedMediaURL1 {
		t.Errorf("Item 1 MediaURL: Expected %s, got %v", expectedMediaURL1, items[0].MediaURL)
	}
	if items[0].ProfileLink != expectedProfileLink1 {
		t.Errorf("Item 1 ProfileLink: Expected %s, got %s", expectedProfileLink1, items[0].ProfileLink)
	}
	// Check timestamp is recent (within a reasonable margin)
	if time.Since(items[0].Timestamp) > 5*time.Hour+5*time.Minute || time.Since(items[0].Timestamp) < 5*time.Hour-5*time.Minute {
		t.Errorf("Item 1 Timestamp is not as expected: %v (expected around %v ago)", items[0].Timestamp, 5*time.Hour)
	}
	if items[0].Interactions != expectedInteractions1 {
		t.Errorf("Item 1 Interactions: Expected %d, got %d", expectedInteractions1, items[0].Interactions)
	}

	// Test second item
	expectedPlatform2 := "instagram"
	expectedPostContent2 := "New recipe alert! Delicious homemade pasta."
	expectedUsername2 := "FoodieChef"
	expectedMediaURL2 := "https://instagram.com/p/pasta.jpg"
	expectedProfileLink2 := "https://instagram.com/foodiechef"
	expectedInteractions2 := 280

	if items[1].Platform != expectedPlatform2 {
		t.Errorf("Item 2 Platform: Expected %s, got %s", expectedPlatform2, items[1].Platform)
	}
	if items[1].PostContent != expectedPostContent2 {
		t.Errorf("Item 2 PostContent: Expected %s, got %s", expectedPostContent2, items[1].PostContent)
	}
	if items[1].Username != expectedUsername2 {
		t.Errorf("Item 2 Username: Expected %s, got %s", expectedUsername2, items[1].Username)
	}
	if items[1].MediaURL == nil || *items[1].MediaURL != expectedMediaURL2 {
		t.Errorf("Item 2 MediaURL: Expected %s, got %v", expectedMediaURL2, items[1].MediaURL)
	}
	if items[1].ProfileLink != expectedProfileLink2 {
		t.Errorf("Item 2 ProfileLink: Expected %s, got %s", expectedProfileLink2, items[1].ProfileLink)
	}
	if time.Since(items[1].Timestamp) > 20*time.Hour+5*time.Minute || time.Since(items[1].Timestamp) < 20*time.Hour-5*time.Minute {
		t.Errorf("Item 2 Timestamp is not as expected: %v (expected around %v ago)", items[1].Timestamp, 20*time.Hour)
	}
	if items[1].Interactions != expectedInteractions2 {
		t.Errorf("Item 2 Interactions: Expected %d, got %d", expectedInteractions2, items[1].Interactions)
	}
}

func TestInstagramFeed_Fetch_MissingAPIKeys(t *testing.T) {
	// Unset environment variables to simulate missing keys
	os.Unsetenv("INSTAGRAM_API_KEY")
	os.Unsetenv("INSTAGRAM_API_SECRET")

	instagramFeed := NewInstagramFeed()
	_, err := instagramFeed.Fetch()

	if err == nil {
		t.Fatalf("Expected an error for missing API keys, got nil")
	}
	expectedError := "Instagram API keys not set in environment variables"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
