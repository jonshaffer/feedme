package rss

import (
	"net/http"
	"net/http/httptest"
	"strings" // Added for strings.Contains
	"testing"
	"time"

	"feed/feeds"
)

func init() {
	// Dummy usage to ensure feeds import is not removed by linter
	_ = feeds.FeedItem{}
}

func TestRSSFeed_Fetch(t *testing.T) {
	// Mock RSS feed content
	mockRSSContent := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test Blog</title>
    <link>http://testblog.com</link>
    <item>
      <title>First Post</title>
      <link>http://testblog.com/first</link>
      <description>This is the first post content.</description>
      <pubDate>Mon, 01 Jan 2025 12:00:00 +0000</pubDate>
      <author>Author One</author>
    </item>
    <item>
      <title>Second Post</title>
      <link>http://testblog.com/second</link>
      <description>This is the second post content.</description>
      <pubDate>Sun, 31 Dec 2024 10:00:00 +0000</pubDate>
    </item>
  </channel>
</rss>`

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(mockRSSContent))
		if err != nil {
			t.Fatalf("Failed to write mock response: %v", err)
		}
	}))
	defer server.Close()

	rssFeed := NewRSSFeed(server.URL)
	items, err := rssFeed.Fetch()

	if err != nil {
		t.Fatalf("Fetch returned an error: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Test first item
	expectedPlatform1 := "rss"
	expectedTitle1 := "First Post"
	expectedDescription1 := "This is the first post content."
	expectedPostContent1 := expectedTitle1 + "\n" + expectedDescription1
	expectedUsername1 := "Author One"
	expectedProfileLink1 := "http://testblog.com/first"
	expectedTimestamp1, _ := time.Parse(time.RFC1123Z, "Mon, 01 Jan 2025 12:00:00 +0000")

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
	if !items[0].Timestamp.Equal(expectedTimestamp1) {
		t.Errorf("Item 1 Timestamp: Expected %v, got %v", expectedTimestamp1, items[0].Timestamp)
	}
	if items[0].Interactions != 0 {
		t.Errorf("Item 1 Interactions: Expected 0, got %d", items[0].Interactions)
	}

	// Test second item (without author)
	expectedPlatform2 := "rss"
	expectedTitle2 := "Second Post"
	expectedDescription2 := "This is the second post content."
	expectedPostContent2 := expectedTitle2 + "\n" + expectedDescription2
	expectedUsername2 := "Test Blog" // Should fall back to channel title
	expectedProfileLink2 := "http://testblog.com/second"
	expectedTimestamp2, _ := time.Parse(time.RFC1123Z, "Sun, 31 Dec 2024 10:00:00 +0000")

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
	if !items[1].Timestamp.Equal(expectedTimestamp2) {
		t.Errorf("Item 2 Timestamp: Expected %v, got %v", expectedTimestamp2, items[1].Timestamp)
	}
	if items[1].Interactions != 0 {
		t.Errorf("Item 2 Interactions: Expected 0, got %d", items[1].Interactions)
	}
}

func TestRSSFeed_Fetch_InvalidURL(t *testing.T) {
	rssFeed := NewRSSFeed("http://invalid-url-that-does-not-exist.com")
	_, err := rssFeed.Fetch()

	if err == nil {
		t.Fatalf("Expected an error for invalid URL, got nil")
	}
}

func TestRSSFeed_Fetch_Non200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // Simulate 404
	}))
	defer server.Close()

	rssFeed := NewRSSFeed(server.URL)
	_, err := rssFeed.Fetch()

	if err == nil {
		t.Fatalf("Expected an error for non-200 status, got nil")
	}
	expectedErrorPart := "status code: 404"
	if !contains(err.Error(), expectedErrorPart) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedErrorPart, err.Error())
	}
}

func TestRSSFeed_Fetch_MalformedXML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("<invalid>xml<feed>")) // Malformed XML
		if err != nil {
			t.Fatalf("Failed to write mock response: %v", err)
		}
	}))
	defer server.Close()

	rssFeed := NewRSSFeed(server.URL)
	_, err := rssFeed.Fetch()

	if err == nil {
		t.Fatalf("Expected an error for malformed XML, got nil")
	}
	expectedErrorPart := "failed to unmarshal RSS feed"
	if !contains(err.Error(), expectedErrorPart) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedErrorPart, err.Error())
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
