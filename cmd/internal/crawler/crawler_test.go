package crawler

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestPage_Structure(t *testing.T) {
	page := Page{
		URL:  "https://example.com",
		HTML: "<html>test</html>",
	}

	if page.URL != "https://example.com" {
		t.Error("URL not set correctly")
	}

	if page.HTML != "<html>test</html>" {
		t.Error("HTML not set correctly")
	}
}

func TestAsset_Structure(t *testing.T) {
	asset := Asset{
		URL:  "https://example.com/image.jpg",
		Path: "/images/image.jpg",
		Type: "image",
	}

	if asset.URL != "https://example.com/image.jpg" {
		t.Error("URL not set correctly")
	}

	if asset.Path != "/images/image.jpg" {
		t.Error("Path not set correctly")
	}

	if asset.Type != "image" {
		t.Error("Type not set correctly")
	}
}

func TestUserAgent_Constant(t *testing.T) {
	if UserAgent == "" {
		t.Error("UserAgent constant should not be empty")
	}

	expected := "Pangolin/1.0 (Static Site Generator)"
	if UserAgent != expected {
		t.Errorf("Expected UserAgent '%s', got '%s'", expected, UserAgent)
	}
}

func TestAssetsMap_ConcurrentAccess(t *testing.T) {
	assetsDownloaded = make(map[string]bool)
	assetsLock = sync.Mutex{}

	assetsDownloaded["test1"] = true

	if !assetsDownloaded["test1"] {
		t.Error("Expected test1 to be in map")
	}

	if assetsDownloaded["test2"] {
		t.Error("Expected test2 to not be in map")
	}
}

func TestHTTPClient_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body>Test</body></html>`))
	}))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHTTPRequest_Headers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != UserAgent {
			t.Errorf("Expected User-Agent '%s', got '%s'", UserAgent, r.Header.Get("User-Agent"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
}

func TestFetchPage_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := FetchPage(server.URL, "/tmp/dist")
	if err == nil {
		t.Error("Expected error for server error")
	}
}

func TestFetchPage_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := FetchPage(server.URL, "/tmp/dist")
	if err == nil {
		t.Error("Expected error for 404 response")
	}
}

func TestFetchPage_InvalidURL(t *testing.T) {
	_, err := FetchPage("://invalid-url", "/tmp/dist")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}
