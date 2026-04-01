package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig_FileNotFound(t *testing.T) {
	viper.Reset()
	viper.SetConfigName("nonexistent")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/tmp/nonexistentpath123")

	_, err := LoadConfig()
	if err == nil {
		t.Error("Expected error when config file not found")
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "pangolin.yaml")

	viper.Reset()
	viper.SetConfigName("pangolin")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(tmpDir)

	cfg := &Config{
		SiteURL: "https://example.com",
		APIKey:  "test-api-key",
	}

	viper.Set("site_url", cfg.SiteURL)
	viper.Set("api_key", cfg.APIKey)

	if err := viper.WriteConfigAs(configPath); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestConfig_Structure(t *testing.T) {
	cfg := &Config{
		SiteURL:      "https://test.com",
		APIKey:       "key123",
		S3Bucket:     "my-bucket",
		S3Region:     "us-west-2",
		NetlifyToken: "netlify-token",
		NetlifySite:  "my-site",
	}

	if cfg.SiteURL != "https://test.com" {
		t.Errorf("Expected SiteURL 'https://test.com', got '%s'", cfg.SiteURL)
	}
	if cfg.APIKey != "key123" {
		t.Errorf("Expected APIKey 'key123', got '%s'", cfg.APIKey)
	}
	if cfg.S3Bucket != "my-bucket" {
		t.Errorf("Expected S3Bucket 'my-bucket', got '%s'", cfg.S3Bucket)
	}
	if cfg.S3Region != "us-west-2" {
		t.Errorf("Expected S3Region 'us-west-2', got '%s'", cfg.S3Region)
	}
	if cfg.NetlifyToken != "netlify-token" {
		t.Errorf("Expected NetlifyToken 'netlify-token', got '%s'", cfg.NetlifyToken)
	}
	if cfg.NetlifySite != "my-site" {
		t.Errorf("Expected NetlifySite 'my-site', got '%s'", cfg.NetlifySite)
	}
}

func TestConfig_EmptyFields(t *testing.T) {
	cfg := &Config{}

	if cfg.SiteURL != "" {
		t.Error("Expected empty SiteURL")
	}
	if cfg.APIKey != "" {
		t.Error("Expected empty APIKey")
	}
	if cfg.S3Bucket != "" {
		t.Error("Expected empty S3Bucket")
	}
}
