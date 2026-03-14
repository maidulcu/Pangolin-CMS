package sitemap

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pangolin-cms/staticpress/cmd/internal/config"

	"github.com/PuerkitoBio/goquery"
)

func FetchSitemaps() ([]string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	sitemapURL := cfg.SiteURL + "/sitemap.xml"
	urls, err := fetchSitemap(sitemapURL)
	if err == nil {
		return urls, nil
	}

	sitemapURL = cfg.SiteURL + "/wp-sitemap.xml"
	return fetchSitemap(sitemapURL)
}

func fetchSitemap(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	cfg, _ := config.LoadConfig()
	if cfg != nil && cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("sitemap not found")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var urls []string

	doc.Find("loc").Each(func(i int, s *goquery.Selection) {
		loc := s.Text()
		if loc != "" && !strings.Contains(loc, "sitemap") {
			urls = append(urls, loc)
		}
	})

	return urls, nil
}
