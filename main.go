package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"feed/config"
	"feed/feeds"
	"feed/feeds/credly"
	"feed/feeds/goodreads"
	"feed/feeds/instagram"
	"feed/feeds/linkedin"
	"feed/feeds/reddit"
	"feed/feeds/rss"
	"feed/feeds/strava"
	"feed/feeds/threads"
	"feed/feeds/x"
)

// PaginatedFeed represents the structure for a paginated JSON output.
type PaginatedFeed struct {
	Items       []feeds.FeedItem `json:"items"`
	CurrentPage int              `json:"current_page"`
	TotalPages  int              `json:"total_pages"`
	NextPage    *string          `json:"next_page,omitempty"`
	PrevPage    *string          `json:"prev_page,omitempty"`
}

// MetaData represents the overall metadata for the generated feeds.
type MetaData struct {
	TotalItems      int               `json:"total_items"`
	TotalPages      int               `json:"total_pages,omitempty"`
	MainFeedPages   []string          `json:"main_feed_pages,omitempty"`
	PlatformFeeds   map[string]string `json:"platform_feeds,omitempty"`
	IndividualItems string            `json:"individual_items_directory,omitempty"`
}

func main() {
	// Dummy usage of time to ensure import is not removed by linter
	_ = time.Now()

	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Printf("Warning: Could not load config.yaml (%v). Attempting to load config.yaml.example.\n", err)
		cfg, err = config.LoadConfig("config/config.yaml.example")
		if err != nil {
			log.Fatalf("Error: Could not load config.yaml or config.yaml.example: %v", err)
		}
		log.Println("Successfully loaded config.yaml.example.")
	}

	var allFeedItems []feeds.FeedItem

	// Initialize and fetch data from enabled feeds
	if cfg.Feeds.LinkedIn.Enabled {
		log.Println("Fetching LinkedIn feed...")
		liFeed := linkedin.NewLinkedInFeed()
		items, err := liFeed.Fetch()
		if err != nil {
			log.Printf("Error fetching LinkedIn feed: %v", err)
		} else {
			allFeedItems = append(allFeedItems, items...)
			log.Printf("Fetched %d items from LinkedIn.", len(items))
		}
	}

	if cfg.Feeds.Threads.Enabled {
		log.Println("Fetching Threads feed...")
		threadsFeed := threads.NewThreadsFeed()
		items, err := threadsFeed.Fetch()
		if err != nil {
			log.Printf("Error fetching Threads feed: %v", err)
		} else {
			allFeedItems = append(allFeedItems, items...)
			log.Printf("Fetched %d items from Threads.", len(items))
		}
	}

	if cfg.Feeds.X.Enabled {
		log.Println("Fetching X feed...")
		xFeed := x.NewXFeed()
		items, err := xFeed.Fetch()
		if err != nil {
			log.Printf("Error fetching X feed: %v", err)
		} else {
			allFeedItems = append(allFeedItems, items...)
			log.Printf("Fetched %d items from X.", len(items))
		}
	}

	if cfg.Feeds.Instagram.Enabled {
		log.Println("Fetching Instagram feed...")
		instagramFeed := instagram.NewInstagramFeed()
		items, err := instagramFeed.Fetch()
		if err != nil {
			log.Printf("Error fetching Instagram feed: %v", err)
		} else {
			allFeedItems = append(allFeedItems, items...)
			log.Printf("Fetched %d items from Instagram.", len(items))
		}
	}

	if cfg.Feeds.Reddit.Enabled {
		log.Println("Fetching Reddit feed...")
		redditFeed := reddit.NewRedditFeed()
		items, err := redditFeed.Fetch()
		if err != nil {
			log.Printf("Error fetching Reddit feed: %v", err)
		} else {
			allFeedItems = append(allFeedItems, items...)
			log.Printf("Fetched %d items from Reddit.", len(items))
		}
	}

	if cfg.Feeds.Strava.Enabled {
		log.Println("Fetching Strava feed...")
		stravaFeed := strava.NewStravaFeed()
		items, err := stravaFeed.Fetch()
		if err != nil {
			log.Printf("Error fetching Strava feed: %v", err)
		} else {
			allFeedItems = append(allFeedItems, items...)
			log.Printf("Fetched %d items from Strava.", len(items))
		}
	}

	if cfg.Feeds.Goodreads.Enabled {
		log.Println("Fetching Goodreads feed...")
		goodreadsFeed := goodreads.NewGoodreadsFeed()
		items, err := goodreadsFeed.Fetch()
		if err != nil {
			log.Printf("Error fetching Goodreads feed: %v", err)
		} else {
			allFeedItems = append(allFeedItems, items...)
			log.Printf("Fetched %d items from Goodreads.", len(items))
		}
	}

	if cfg.Feeds.Credly.Enabled {
		log.Println("Fetching Credly feed...")
		credlyFeed := credly.NewCredlyFeed()
		items, err := credlyFeed.Fetch()
		if err != nil {
			log.Printf("Error fetching Credly feed: %v", err)
		} else {
			allFeedItems = append(allFeedItems, items...)
			log.Printf("Fetched %d items from Credly.", len(items))
		}
	}

	if cfg.Feeds.RSS.Enabled {
		log.Println("Fetching RSS feeds...")
		for _, url := range cfg.Feeds.RSS.URLs {
			rssFeed := rss.NewRSSFeed(url)
			items, err := rssFeed.Fetch()
			if err != nil {
				log.Printf("Error fetching RSS feed from %s: %v", url, err)
			} else {
				allFeedItems = append(allFeedItems, items...)
				log.Printf("Fetched %d items from RSS feed: %s.", len(items), url)
			}
		}
	}

	// Sort feed items by timestamp in descending order
	sort.Slice(allFeedItems, func(i, j int) bool {
		return allFeedItems[i].Timestamp.After(allFeedItems[j].Timestamp)
	})

	// Apply output limit if configured
	if cfg.OutputLimit > 0 && len(allFeedItems) > cfg.OutputLimit {
		log.Printf("Applying output limit: truncating %d items to %d.", len(allFeedItems), cfg.OutputLimit)
		allFeedItems = allFeedItems[:cfg.OutputLimit]
	}

	// Create output directories if they don't exist
	outputDir := "output"
	itemsDir := filepath.Join(outputDir, "items")
	platformsDir := filepath.Join(outputDir, "platforms")

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			log.Fatalf("Error creating output directory: %v", err)
		}
	}

	var meta MetaData
	meta.TotalItems = len(allFeedItems)
	meta.PlatformFeeds = make(map[string]string)

	// Generate individual item files if enabled
	if cfg.GenerateIndividualItemFiles {
		log.Println("Generating individual item files...")
		if _, err := os.Stat(itemsDir); os.IsNotExist(err) {
			err = os.MkdirAll(itemsDir, 0755)
			if err != nil {
				log.Fatalf("Error creating individual items directory: %v", err)
			}
		}
		meta.IndividualItems = "items/"

		for i, item := range allFeedItems {
			// Use a combination of timestamp and index for a unique filename
			filename := fmt.Sprintf("%s_%d.json", item.Timestamp.Format("20060102150405"), i)
			itemFilePath := filepath.Join(itemsDir, filename)
			item.Permalink = filepath.Join("items", filename) // Set permalink for the item

			jsonData, err := json.MarshalIndent(item, "", "  ")
			if err != nil {
				log.Printf("Error marshalling individual item JSON for %s: %v", filename, err)
				continue
			}

			err = ioutil.WriteFile(itemFilePath, jsonData, 0644)
			if err != nil {
				log.Printf("Error writing individual item JSON to file %s: %v", itemFilePath, err)
			}
			allFeedItems[i] = item // Update the item in the slice with the permalink
		}
		log.Printf("Generated %d individual item files.", len(allFeedItems))
	}

	// Generate platform-specific feeds if enabled
	if cfg.GeneratePlatformFeeds {
		log.Println("Generating platform-specific feeds...")
		if _, err := os.Stat(platformsDir); os.IsNotExist(err) {
			err = os.MkdirAll(platformsDir, 0755)
			if err != nil {
				log.Fatalf("Error creating platforms directory: %v", err)
			}
		}

		platformItems := make(map[string][]feeds.FeedItem)
		for _, item := range allFeedItems {
			platformItems[item.Platform] = append(platformItems[item.Platform], item)
		}

		for platform, items := range platformItems {
			log.Printf("Generating feed for platform: %s with %d items...", platform, len(items))
			platformSlug := strings.ToLower(platform) // Create a URL-friendly slug

			if cfg.PageSize > 1 {
				// Paginate platform feed
				totalPages := (len(items) + cfg.PageSize - 1) / cfg.PageSize
				for pageNum := 0; pageNum < totalPages; pageNum++ {
					start := pageNum * cfg.PageSize
					end := (pageNum + 1) * cfg.PageSize
					if end > len(items) {
						end = len(items)
					}
					pageItems := items[start:end]

					nextPage := (*string)(nil)
					if pageNum < totalPages-1 {
						nextPageFile := fmt.Sprintf("%s_page_%d.json", platformSlug, pageNum+2)
						nextPage = &nextPageFile
					}

					prevPage := (*string)(nil)
					if pageNum > 0 {
						prevPageFile := fmt.Sprintf("%s_page_%d.json", platformSlug, pageNum)
						prevPage = &prevPageFile
					}

					paginatedFeed := PaginatedFeed{
						Items:       pageItems,
						CurrentPage: pageNum + 1,
						TotalPages:  totalPages,
						NextPage:    nextPage,
						PrevPage:    prevPage,
					}

					filename := fmt.Sprintf("%s_page_%d.json", platformSlug, pageNum+1)
					filePath := filepath.Join(platformsDir, filename)
					jsonData, err := json.MarshalIndent(paginatedFeed, "", "  ")
					if err != nil {
						log.Printf("Error marshalling JSON for platform %s page %d: %v", platform, pageNum+1, err)
						continue
					}
					err = ioutil.WriteFile(filePath, jsonData, 0644)
					if err != nil {
						log.Printf("Error writing JSON for platform %s page %d to file: %v", platform, pageNum+1, err)
					}
					if pageNum == 0 { // Store link to the first page of each platform feed
						meta.PlatformFeeds[platform] = filepath.Join("platforms", filename)
					}
				}
			} else {
				// Single file for platform feed
				filename := fmt.Sprintf("%s.json", platformSlug)
				filePath := filepath.Join(platformsDir, filename)
				jsonData, err := json.MarshalIndent(items, "", "  ")
				if err != nil {
					log.Printf("Error marshalling JSON for platform %s: %v", platform, err)
					continue
				}
				err = ioutil.WriteFile(filePath, jsonData, 0644)
				if err != nil {
					log.Printf("Error writing JSON for platform %s to file: %v", platform, err)
				}
				meta.PlatformFeeds[platform] = filepath.Join("platforms", filename)
			}
		}
		log.Println("Finished generating platform-specific feeds.")
	}

	// Generate main aggregated feed (with pagination if configured)
	log.Println("Generating main aggregated feed...")
	if cfg.PageSize > 1 {
		totalPages := (len(allFeedItems) + cfg.PageSize - 1) / cfg.PageSize
		meta.TotalPages = totalPages
		for pageNum := 0; pageNum < totalPages; pageNum++ {
			start := pageNum * cfg.PageSize
			end := (pageNum + 1) * cfg.PageSize
			if end > len(allFeedItems) {
				end = len(allFeedItems)
			}
			pageItems := allFeedItems[start:end]

			nextPage := (*string)(nil)
			if pageNum < totalPages-1 {
				nextPageFile := fmt.Sprintf("feed_page_%d.json", pageNum+2)
				nextPage = &nextPageFile
			}

			prevPage := (*string)(nil)
			if pageNum > 0 {
				prevPageFile := fmt.Sprintf("feed_page_%d.json", pageNum)
				prevPage = &prevPageFile
			}

			paginatedFeed := PaginatedFeed{
				Items:       pageItems,
				CurrentPage: pageNum + 1,
				TotalPages:  totalPages,
				NextPage:    nextPage,
				PrevPage:    prevPage,
			}

			filename := fmt.Sprintf("feed_page_%d.json", pageNum+1)
			outputFilePath := filepath.Join(outputDir, filename)
			jsonData, err := json.MarshalIndent(paginatedFeed, "", "  ")
			if err != nil {
				log.Fatalf("Error marshalling JSON for main feed page %d: %v", pageNum+1, err)
			}
			err = ioutil.WriteFile(outputFilePath, jsonData, 0644)
			if err != nil {
				log.Fatalf("Error writing JSON for main feed page %d to file: %v", pageNum+1, err)
			}
			meta.MainFeedPages = append(meta.MainFeedPages, filepath.Join(outputDir, filename))
		}
		log.Printf("Successfully aggregated %d feed items into %d paginated files.", len(allFeedItems), totalPages)
	} else {
		// Write aggregated data to a single JSON file
		outputFilePath := filepath.Join(outputDir, "feed.json")
		jsonData, err := json.MarshalIndent(allFeedItems, "", "  ")
		if err != nil {
			log.Fatalf("Error marshalling JSON: %v", err)
		}

		err = ioutil.WriteFile(outputFilePath, jsonData, 0644)
		if err != nil {
			log.Fatalf("Error writing JSON to file: %v", err)
		}
		meta.MainFeedPages = []string{filepath.Join(outputDir, "feed.json")}
		log.Printf("Successfully aggregated %d feed items to %s", len(allFeedItems), outputFilePath)
	}

	// Write metadata to meta.json
	metaFilePath := filepath.Join(outputDir, "meta.json")
	metaJsonData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling metadata JSON: %v", err)
	}
	err = ioutil.WriteFile(metaFilePath, metaJsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing metadata JSON to file: %v", err)
	}
	log.Printf("Successfully generated metadata to %s", metaFilePath)
}
