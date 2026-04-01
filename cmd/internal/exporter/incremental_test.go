package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewIncrementalExporter(t *testing.T) {
	tmpDir := t.TempDir()

	exporter := NewIncrementalExporter(tmpDir, 5)
	if exporter == nil {
		t.Fatal("Expected non-nil IncrementalExporter")
	}

	if exporter.cache == nil {
		t.Error("Expected non-nil cache")
	}

	if exporter.cache.Pages == nil {
		t.Error("Expected non-nil cache.Pages")
	}

	expectedCacheDir := filepath.Join(tmpDir, ".pangolin")
	if _, err := os.Stat(expectedCacheDir); os.IsNotExist(err) {
		t.Error("Expected .pangolin directory to be created")
	}
}

func TestLoadCache(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "export_cache.json")

	cacheData := ExportCache{
		Pages: map[string]PageCache{
			"https://example.com/page1": {
				ETag:        "abc123",
				ContentHash: "hash123",
				ExportedAt:  time.Now().Format(time.RFC3339),
			},
		},
	}

	data, _ := json.Marshal(cacheData)
	os.WriteFile(cachePath, data, 0644)

	cache := loadCache(cachePath)
	if cache == nil {
		t.Fatal("Expected non-nil cache")
	}

	if len(cache.Pages) != 1 {
		t.Errorf("Expected 1 cached page, got %d", len(cache.Pages))
	}

	if _, exists := cache.Pages["https://example.com/page1"]; !exists {
		t.Error("Expected page1 to be in cache")
	}
}

func TestLoadCache_FileNotFound(t *testing.T) {
	cache := loadCache("/nonexistent/path/cache.json")
	if cache == nil {
		t.Fatal("Expected non-nil cache for non-existent file")
	}

	if len(cache.Pages) != 0 {
		t.Errorf("Expected empty cache, got %d pages", len(cache.Pages))
	}
}

func TestLoadCache_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "export_cache.json")

	os.WriteFile(cachePath, []byte("invalid json"), 0644)

	cache := loadCache(cachePath)
	if cache == nil {
		t.Fatal("Expected non-nil cache for invalid JSON")
	}

	if len(cache.Pages) != 0 {
		t.Errorf("Expected empty cache, got %d pages", len(cache.Pages))
	}
}

func TestExportCache_Structure(t *testing.T) {
	cache := &ExportCache{
		Pages: map[string]PageCache{
			"https://example.com": {
				ETag:         "etag123",
				LastModified: "2024-01-01",
				ContentHash:  "hash123",
				ExportedAt:   "2024-01-01T00:00:00Z",
			},
		},
	}

	data, err := json.Marshal(cache)
	if err != nil {
		t.Fatalf("Failed to marshal cache: %v", err)
	}

	var unmarshaled ExportCache
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal cache: %v", err)
	}

	if unmarshaled.Pages["https://example.com"].ETag != "etag123" {
		t.Error("Expected ETag to be preserved")
	}
}

func TestPageCache_Structure(t *testing.T) {
	pc := PageCache{
		ETag:         "test-etag",
		LastModified: "test-date",
		ContentHash:  "test-hash",
		ExportedAt:   "2024-01-01T00:00:00Z",
	}

	if pc.ETag != "test-etag" {
		t.Error("ETag not set correctly")
	}
	if pc.LastModified != "test-date" {
		t.Error("LastModified not set correctly")
	}
	if pc.ContentHash != "test-hash" {
		t.Error("ContentHash not set correctly")
	}
}

func TestIncrementalExporter_ClearCache(t *testing.T) {
	tmpDir := t.TempDir()

	exporter := NewIncrementalExporter(tmpDir, 5)
	exporter.cache.Pages["https://example.com"] = PageCache{
		ETag:       "test",
		ExportedAt: time.Now().Format(time.RFC3339),
	}

	exporter.saveCache()

	err := exporter.ClearCache()
	if err != nil {
		t.Fatalf("ClearCache failed: %v", err)
	}

	if len(exporter.cache.Pages) != 0 {
		t.Errorf("Expected 0 pages after clear, got %d", len(exporter.cache.Pages))
	}
}

func TestIncrementalExporter_GetCacheStats(t *testing.T) {
	tmpDir := t.TempDir()

	exporter := NewIncrementalExporter(tmpDir, 5)
	exporter.cache.Pages["page1"] = PageCache{ExportedAt: "2024-01-01T00:00:00Z"}
	exporter.cache.Pages["page2"] = PageCache{ExportedAt: "2024-01-02T00:00:00Z"}
	exporter.cache.Pages["page3"] = PageCache{ExportedAt: "2024-01-03T00:00:00Z"}

	count, oldest := exporter.GetCacheStats()

	if count != 3 {
		t.Errorf("Expected 3 pages, got %d", count)
	}

	if oldest.IsZero() {
		t.Error("Expected non-zero oldest time")
	}
}

func TestIncrementalExporter_GetCacheStats_Empty(t *testing.T) {
	tmpDir := t.TempDir()

	exporter := NewIncrementalExporter(tmpDir, 5)

	count, oldest := exporter.GetCacheStats()

	if count != 0 {
		t.Errorf("Expected 0 pages, got %d", count)
	}

	if !oldest.IsZero() {
		t.Error("Expected zero time for empty cache")
	}
}

func TestIncrementalExporter_RemoveFromCache(t *testing.T) {
	tmpDir := t.TempDir()

	exporter := NewIncrementalExporter(tmpDir, 5)
	exporter.cache.Pages["page1"] = PageCache{ETag: "test1"}
	exporter.cache.Pages["page2"] = PageCache{ETag: "test2"}

	exporter.RemoveFromCache("page1")

	if len(exporter.cache.Pages) != 1 {
		t.Errorf("Expected 1 page after removal, got %d", len(exporter.cache.Pages))
	}

	if _, exists := exporter.cache.Pages["page1"]; exists {
		t.Error("page1 should have been removed")
	}

	if _, exists := exporter.cache.Pages["page2"]; !exists {
		t.Error("page2 should still exist")
	}
}

func TestSaveCache(t *testing.T) {
	tmpDir := t.TempDir()

	exporter := NewIncrementalExporter(tmpDir, 5)
	exporter.cache.Pages["https://example.com"] = PageCache{
		ETag:       "test-etag",
		ExportedAt: time.Now().Format(time.RFC3339),
	}

	err := exporter.saveCache()
	if err != nil {
		t.Fatalf("saveCache failed: %v", err)
	}

	cacheFile := filepath.Join(tmpDir, ".pangolin", "export_cache.json")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Error("Cache file should exist after saveCache")
	}
}
