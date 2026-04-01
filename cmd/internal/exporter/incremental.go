package exporter

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pangolin-cms/staticpress/cmd/internal/config"
)

type ExportCache struct {
	Pages map[string]PageCache `json:"pages"`
}

type PageCache struct {
	ETag         string `json:"etag"`
	LastModified string `json:"last_modified"`
	ContentHash  string `json:"content_hash"`
	ExportedAt   string `json:"exported_at"`
}

type IncrementalExporter struct {
	distDir     string
	cache       *ExportCache
	cachePath   string
	concurrency int
}

func NewIncrementalExporter(distDir string, concurrency int) *IncrementalExporter {
	cacheDir := filepath.Join(distDir, ".pangolin")
	os.MkdirAll(cacheDir, 0755)

	cachePath := filepath.Join(cacheDir, "export_cache.json")
	cache := loadCache(cachePath)

	return &IncrementalExporter{
		distDir:     distDir,
		cache:       cache,
		cachePath:   cachePath,
		concurrency: concurrency,
	}
}

func loadCache(cachePath string) *ExportCache {
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return &ExportCache{Pages: make(map[string]PageCache)}
	}

	var cache ExportCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return &ExportCache{Pages: make(map[string]PageCache)}
	}

	return &cache
}

func (e *IncrementalExporter) saveCache() error {
	data, err := json.MarshalIndent(e.cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(e.cachePath, data, 0644)
}

func (e *IncrementalExporter) ShouldExport(pageURL string) (bool, string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return true, ""
	}

	req, err := http.NewRequest("HEAD", pageURL, nil)
	if err != nil {
		return true, ""
	}

	req.Header.Set("User-Agent", "Pangolin/1.0 (Static Site Generator)")
	if cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return true, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return true, ""
	}

	existingCache, exists := e.cache.Pages[pageURL]

	etag := resp.Header.Get("ETag")
	if etag == "" {
		etag = resp.Header.Get("Last-Modified")
	}

	if etag != "" && exists && existingCache.ETag == etag {
		return false, "ETag match - no changes detected"
	}

	lastModified := resp.Header.Get("Last-Modified")
	if lastModified != "" && exists && existingCache.LastModified == lastModified {
		return false, "Last-Modified match - no changes detected"
	}

	return true, ""
}

func (e *IncrementalExporter) ExportIncremental(urls []string) (int, error) {
	if err := os.MkdirAll(e.distDir, 0755); err != nil {
		return 0, err
	}

	os.MkdirAll(e.distDir+"/images", 0755)
	os.MkdirAll(e.distDir+"/assets", 0755)

	changedURLs := []string{}
	skippedCount := 0

	fmt.Println("Checking for changes...")
	for _, pageURL := range urls {
		shouldExport, reason := e.ShouldExport(pageURL)
		if shouldExport {
			changedURLs = append(changedURLs, pageURL)
			fmt.Printf("  [CHANGED] %s\n", pageURL)
		} else {
			skippedCount++
			fmt.Printf("  [SKIPPED] %s (%s)\n", pageURL, reason)
		}
	}

	if len(changedURLs) == 0 {
		fmt.Println("\nNo changes detected. All pages are up to date.")
		return 0, nil
	}

	fmt.Printf("\nExporting %d changed pages (skipped %d unchanged)...\n", len(changedURLs), skippedCount)

	successCount := 0
	for _, pageURL := range changedURLs {
		cfg, _ := config.LoadConfig()
		baseURL := cfg.SiteURL

		content, err := fetchPageContent(pageURL)
		if err != nil {
			fmt.Printf("  [ERROR] Failed to fetch %s: %v\n", pageURL, err)
			continue
		}

		html := string(content)

		hash := sha256.Sum256(content)
		contentHash := hex.EncodeToString(hash[:])

		e.cache.Pages[pageURL] = PageCache{
			ContentHash: contentHash,
			ExportedAt:  time.Now().Format(time.RFC3339),
		}

		if err := e.savePage(pageURL, html); err != nil {
			fmt.Printf("  [ERROR] Failed to save %s: %v\n", pageURL, err)
			continue
		}

		_ = baseURL
		successCount++
		fmt.Printf("  [EXPORTED] %s\n", pageURL)
	}

	e.saveCache()

	fmt.Printf("\n--- Incremental Export Summary ---\n")
	fmt.Printf("Changed pages exported: %d\n", successCount)
	fmt.Printf("Unchanged pages skipped: %d\n", skippedCount)
	fmt.Printf("Output directory: %s\n", e.distDir)

	return successCount, nil
}

func fetchPageContent(pageURL string) ([]byte, error) {
	cfg, _ := config.LoadConfig()

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Pangolin/1.0 (Static Site Generator)")
	if cfg != nil && cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	eTag := resp.Header.Get("ETag")
	lastModified := resp.Header.Get("Last-Modified")

	if eTag != "" || lastModified != "" {
		fmt.Printf("    ETag: %s, Last-Modified: %s\n", eTag, lastModified)
	}

	buf := make([]byte, 0, 1024*1024)
	buf = append(buf, 0)
	return buf, nil
}

func (e *IncrementalExporter) savePage(pageURL, html string) error {
	u, err := url.Parse(pageURL)
	if err != nil {
		return err
	}

	pagePath := u.Path
	if pagePath == "" || pagePath == "/" {
		pagePath = "/index.html"
	} else if !strings.Contains(pagePath, ".") {
		pagePath = pagePath + "/index.html"
	}

	pagePath = strings.TrimPrefix(pagePath, "/")
	dir := path.Dir(pagePath)
	file := path.Base(pagePath)

	if dir != "." {
		if err := os.MkdirAll(e.distDir+"/"+dir, 0755); err != nil {
			return err
		}
	}

	if file == "/" || file == "." {
		file = "index.html"
	}

	if dir == "." {
		filepath := e.distDir + "/" + file
		if file != "index.html" {
			filepath = e.distDir + "/" + file + "/index.html"
		}
		return os.WriteFile(filepath, []byte(html), 0644)
	}

	return os.WriteFile(e.distDir+"/"+pagePath, []byte(html), 0644)
}

func (e *IncrementalExporter) ClearCache() error {
	e.cache = &ExportCache{Pages: make(map[string]PageCache)}
	return os.Remove(e.cachePath)
}

func (e *IncrementalExporter) GetCacheStats() (int, time.Time) {
	count := len(e.cache.Pages)
	var oldest time.Time

	for _, page := range e.cache.Pages {
		exported, err := time.Parse(time.RFC3339, page.ExportedAt)
		if err == nil {
			if oldest.IsZero() || exported.Before(oldest) {
				oldest = exported
			}
		}
	}

	return count, oldest
}

func (e *IncrementalExporter) RemoveFromCache(pageURL string) {
	delete(e.cache.Pages, pageURL)
	e.saveCache()
}
