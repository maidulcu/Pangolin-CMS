package sitemap

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

type sitemapURL struct {
	Loc string `xml:"loc"`
}

type urlset struct {
	URLs []sitemapURL `xml:"url"`
}

type sitemapindex struct {
	Sitemaps []sitemapURL `xml:"sitemap"`
}

func parseSitemap(data []byte) ([]string, error) {
	var u urlset
	if err := xml.Unmarshal(data, &u); err != nil {
		return nil, err
	}

	urls := make([]string, len(u.URLs))
	for i, url := range u.URLs {
		urls[i] = url.Loc
	}
	return urls, nil
}

func parseSitemapIndex(data []byte) ([]string, error) {
	var si sitemapindex
	if err := xml.Unmarshal(data, &si); err != nil {
		return nil, err
	}

	urls := make([]string, len(si.Sitemaps))
	for i, sm := range si.Sitemaps {
		urls[i] = sm.Loc
	}
	return urls, nil
}

func TestParseSitemap_XML(t *testing.T) {
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url>
		<loc>https://example.com/page1</loc>
		<lastmod>2024-01-01</lastmod>
	</url>
	<url>
		<loc>https://example.com/page2</loc>
		<lastmod>2024-01-02</lastmod>
	</url>
</urlset>`

	urls, err := parseSitemap([]byte(xmlContent))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(urls) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(urls))
	}

	if urls[0] != "https://example.com/page1" {
		t.Errorf("Expected first URL 'https://example.com/page1', got '%s'", urls[0])
	}
	if urls[1] != "https://example.com/page2" {
		t.Errorf("Expected second URL 'https://example.com/page2', got '%s'", urls[1])
	}
}

func TestParseSitemap_WPSitemap(t *testing.T) {
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<sitemap>
		<loc>https://example.com/wp-sitemap-posts-post-1.xml</loc>
	</sitemap>
	<sitemap>
		<loc>https://example.com/wp-sitemap-posts-page-1.xml</loc>
	</sitemap>
</sitemapindex>`

	urls, err := parseSitemapIndex([]byte(xmlContent))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(urls) != 2 {
		t.Errorf("Expected 2 sitemap URLs, got %d", len(urls))
	}
}

func TestParseSitemap_InvalidXML(t *testing.T) {
	invalidXML := `<not-valid-xml`

	_, err := parseSitemap([]byte(invalidXML))
	if err == nil {
		t.Error("Expected error for invalid XML")
	}
}

func TestParseSitemap_EmptyXML(t *testing.T) {
	emptyXML := `<?xml version="1.0"?><urlset></urlset>`

	urls, err := parseSitemap([]byte(emptyXML))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(urls) != 0 {
		t.Errorf("Expected 0 URLs, got %d", len(urls))
	}
}

func TestFetchSitemap_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := fetchSitemap(server.URL)
	if err == nil {
		t.Error("Expected error for server error")
	}
}

func TestFetchSitemap_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchSitemap(server.URL)
	if err == nil {
		t.Error("Expected error for 404 response")
	}
}
