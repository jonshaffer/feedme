package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the overall structure of the configuration file.
type Config struct {
	Feeds                       FeedConfig `yaml:"feeds"`
	OutputLimit                 int        `yaml:"output_limit"`
	PageSize                    int        `yaml:"page_size"`
	GenerateIndividualItemFiles bool       `yaml:"generate_individual_item_files"`
	GeneratePlatformFeeds       bool       `yaml:"generate_platform_feeds"`
}

// FeedConfig defines which social media feeds are enabled.
type FeedConfig struct {
	LinkedIn  PlatformConfig `yaml:"linkedin"`
	Threads   PlatformConfig `yaml:"threads"`
	X         PlatformConfig `yaml:"x"`
	Instagram PlatformConfig `yaml:"instagram"`
	Reddit    PlatformConfig `yaml:"reddit"`
	Strava    PlatformConfig `yaml:"strava"`
	Goodreads PlatformConfig `yaml:"goodreads"`
	Credly    PlatformConfig `yaml:"credly"`
	RSS       RSSConfig      `yaml:"rss"`
}

// PlatformConfig is a generic configuration for a social media platform.
type PlatformConfig struct {
	Enabled bool `yaml:"enabled"`
}

// RSSConfig holds configuration specific to RSS feeds.
type RSSConfig struct {
	Enabled bool     `yaml:"enabled"`
	URLs    []string `yaml:"urls"`
}

// LoadConfig reads the configuration from the specified YAML file.
func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
