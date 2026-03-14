package crawler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pangolin-cms/staticpress/cmd/internal/config"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL     string
	HTML    string
	Headers http.Header
}

func FetchPage(url string) (*Page, error) {
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
		return nil, fmt.Errorf("failed to fetch %s: status %d", url, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	cfg, _ = config.LoadConfig()
	baseURL := cfg.SiteURL

	rewriteLinks(doc, baseURL)

	html, err := doc.Html()
	if err != nil {
		return nil, err
	}

	return &Page{
		URL:     url,
		HTML:    html,
		Headers: resp.Header,
	}, nil
}

func rewriteLinks(doc *goquery.Document, baseURL string) {
	sel := doc.Find("a[href]")
	sel.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		if strings.HasPrefix(href, baseURL) {
			newHref := strings.TrimPrefix(href, baseURL)
			if newHref == "" || newHref == "/" {
				newHref = "/"
			}
			s.SetAttr("href", newHref)
		}
	})

	sel = doc.Find("img[src]")
	sel.Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		if strings.HasPrefix(src, baseURL) {
			newSrc := strings.TrimPrefix(src, baseURL)
			s.SetAttr("src", newSrc)
		}
	})

	sel = doc.Find("link[href]")
	sel.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		if strings.HasPrefix(href, baseURL) {
			newHref := strings.TrimPrefix(href, baseURL)
			s.SetAttr("href", newHref)
		}
	})

	sel = doc.Find("script[src]")
	sel.Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		if strings.HasPrefix(src, baseURL) {
			newSrc := strings.TrimPrefix(src, baseURL)
			s.SetAttr("src", newSrc)
		}
	})
}
