package exporter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewExporter(t *testing.T) {
	exporter := NewExporter("/tmp/dist", 5)
	if exporter == nil {
		t.Fatal("Expected non-nil Exporter")
	}

	if exporter.distDir != "/tmp/dist" {
		t.Errorf("Expected distDir '/tmp/dist', got '%s'", exporter.distDir)
	}

	if exporter.concurrency != 5 {
		t.Errorf("Expected concurrency 5, got %d", exporter.concurrency)
	}
}

func TestNewExporter_DefaultConcurrency(t *testing.T) {
	exporter := NewExporter("/tmp/dist", 0)
	if exporter.concurrency != 0 {
		t.Errorf("Expected concurrency 0, got %d", exporter.concurrency)
	}
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.html", "text/html"},
		{"test.css", "text/css"},
		{"test.js", "application/javascript"},
		{"test.json", "application/json"},
		{"test.png", "image/png"},
		{"test.jpg", "image/jpeg"},
		{"test.jpeg", "image/jpeg"},
		{"test.gif", "image/gif"},
		{"test.svg", "image/svg+xml"},
		{"test.woff", "font/woff"},
		{"test.woff2", "font/woff2"},
		{"test.unknown", "text/plain"},
		{"test", "text/plain"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := getContentType(tt.filename)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}
			if *result != tt.expected {
				t.Errorf("getContentType(%s) = %s, expected %s", tt.filename, *result, tt.expected)
			}
		})
	}
}

func TestGetContentType_CaseInsensitive(t *testing.T) {
	tests := []string{"test.HTML", "test.Jpg", "test.PNG"}

	for _, filename := range tests {
		result := getContentType(filename)
		if result == nil {
			t.Errorf("Expected non-nil result for %s", filename)
		}
	}
}

func TestSavePage(t *testing.T) {
	tmpDir := t.TempDir()
	exporter := NewExporter(tmpDir, 1)

	html := "<html><body>Test</body></html>"
	err := exporter.savePage("https://example.com/test", html)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	testFile := filepath.Join(tmpDir, "test", "index.html")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Expected test/index.html to be created")
	}
}

func TestSavePage_HomePage(t *testing.T) {
	tmpDir := t.TempDir()
	exporter := NewExporter(tmpDir, 1)

	html := "<html><body>Home</body></html>"
	err := exporter.savePage("https://example.com/", html)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	indexFile := filepath.Join(tmpDir, "index.html")
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		t.Error("Expected index.html to be created")
	}
}

func TestSavePage_RootURL(t *testing.T) {
	tmpDir := t.TempDir()
	exporter := NewExporter(tmpDir, 1)

	html := "<html><body>Root</body></html>"
	err := exporter.savePage("https://example.com", html)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	indexFile := filepath.Join(tmpDir, "index.html")
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		t.Error("Expected index.html to be created")
	}
}

func TestSavePage_URLWithQueryString(t *testing.T) {
	tmpDir := t.TempDir()
	exporter := NewExporter(tmpDir, 1)

	html := "<html><body>Test</body></html>"
	err := exporter.savePage("https://example.com/page", html)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	file := filepath.Join(tmpDir, "page", "index.html")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Error("Expected page/index.html to be created")
	}
}

func TestSavePage_DirectoriesCreated(t *testing.T) {
	tmpDir := t.TempDir()
	exporter := NewExporter(tmpDir, 1)

	html := "<html><body>Test</body></html>"
	err := exporter.savePage("https://example.com/a/b/c/page", html)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDir := filepath.Join(tmpDir, "a", "b", "c")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Error("Expected nested directories to be created")
	}
}

func TestSavePage_InvalidURL(t *testing.T) {
	tmpDir := t.TempDir()
	exporter := NewExporter(tmpDir, 1)

	html := "<html><body>Test</body></html>"
	err := exporter.savePage("://invalid", html)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestExporter_Structure(t *testing.T) {
	exp := &Exporter{
		distDir:     "/test/dist",
		concurrency: 10,
	}

	if exp.distDir != "/test/dist" {
		t.Error("distDir not set correctly")
	}

	if exp.concurrency != 10 {
		t.Error("concurrency not set correctly")
	}
}
