package exporter

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/pangolin-cms/staticpress/cmd/internal/config"
	"github.com/pangolin-cms/staticpress/cmd/internal/crawler"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Exporter struct {
	distDir     string
	concurrency int
}

func NewExporter(distDir string, concurrency int) *Exporter {
	return &Exporter{
		distDir:     distDir,
		concurrency: concurrency,
	}
}

func (e *Exporter) Export(urls []string) error {
	if err := os.MkdirAll(e.distDir, 0755); err != nil {
		return err
	}

	sem := make(chan struct{}, e.concurrency)
	var wg sync.WaitGroup
	errChan := make(chan error, len(urls))

	for _, pageURL := range urls {
		wg.Add(1)
		sem <- struct{}{}

		go func(url string) {
			defer wg.Done()
			defer func() { <-sem }()

			page, err := crawler.FetchPage(url)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch %s: %w", url, err)
				return
			}

			if err := e.savePage(page.URL, page.HTML); err != nil {
				errChan <- fmt.Errorf("failed to save %s: %w", url, err)
			}
		}(pageURL)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return nil
}

func (e *Exporter) savePage(pageURL, html string) error {
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

func DeployToS3(distDir, bucket, region string, cfg *config.Config) error {
	ctx := context.Background()

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	entries, err := os.ReadDir(distDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		content, err := os.ReadFile(distDir + "/" + name)
		if err != nil {
			return err
		}

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      &bucket,
			Key:         &name,
			Body:        strings.NewReader(string(content)),
			ContentType: getContentType(name),
		})
		if err != nil {
			return fmt.Errorf("failed to upload %s: %w", name, err)
		}

		fmt.Printf("Uploaded: %s\n", name)
	}

	return nil
}

func getContentType(filename string) *string {
	ext := strings.ToLower(path.Ext(filename))
	contentTypes := map[string]string{
		".html":  "text/html",
		".css":   "text/css",
		".js":    "application/javascript",
		".json":  "application/json",
		".png":   "image/png",
		".jpg":   "image/jpeg",
		".jpeg":  "image/jpeg",
		".gif":   "image/gif",
		".svg":   "image/svg+xml",
		".woff":  "font/woff",
		".woff2": "font/woff2",
	}

	if ct, ok := contentTypes[ext]; ok {
		return &ct
	}
	defaultType := "text/plain"
	return &defaultType
}
