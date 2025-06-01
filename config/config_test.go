package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempConfigFile := "test_config.yaml"
	content := `
feeds:
  linkedin:
    enabled: true
  threads:
    enabled: false
  rss:
    enabled: true
    urls:
      - "http://example.com/feed1.xml"
      - "http://example.com/feed2.xml"
`
	err := ioutil.WriteFile(tempConfigFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	defer os.Remove(tempConfigFile) // Clean up the temporary file

	cfg, err := LoadConfig(tempConfigFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if !cfg.Feeds.LinkedIn.Enabled {
		t.Errorf("Expected LinkedIn to be enabled, got false")
	}
	if cfg.Feeds.Threads.Enabled {
		t.Errorf("Expected Threads to be disabled, got true")
	}
	if !cfg.Feeds.RSS.Enabled {
		t.Errorf("Expected RSS to be enabled, got false")
	}
	if len(cfg.Feeds.RSS.URLs) != 2 {
		t.Errorf("Expected 2 RSS URLs, got %d", len(cfg.Feeds.RSS.URLs))
	}
	if cfg.Feeds.RSS.URLs[0] != "http://example.com/feed1.xml" {
		t.Errorf("Expected first RSS URL to be 'http://example.com/feed1.xml', got '%s'", cfg.Feeds.RSS.URLs[0])
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("non_existent_file.yaml")
	if err == nil {
		t.Errorf("Expected an error for non-existent file, got nil")
	}
}

func TestLoadConfig_InvalidYaml(t *testing.T) {
	tempConfigFile := "invalid_config.yaml"
	content := `
feeds:
  linkedin:
    enabled: true
  invalid: [
` // Malformed YAML
	err := ioutil.WriteFile(tempConfigFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary invalid config file: %v", err)
	}
	defer os.Remove(tempConfigFile)

	_, err = LoadConfig(tempConfigFile)
	if err == nil {
		t.Errorf("Expected an error for invalid YAML, got nil")
	}
}
