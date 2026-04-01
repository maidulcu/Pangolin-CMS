package exporter

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Bundler struct {
	minify      bool
	parallelism int
}

type BundlerOptions struct {
	Minify      bool
	Parallelism int
}

func NewBundler(opts BundlerOptions) *Bundler {
	parallelism := opts.Parallelism
	if parallelism <= 0 {
		parallelism = 4
	}

	return &Bundler{
		minify:      opts.Minify,
		parallelism: parallelism,
	}
}

func (b *Bundler) BundleDirectory(dir string) error {
	assetsDir := filepath.Join(dir, "assets")
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		fmt.Println("No assets directory found, skipping bundling")
		return nil
	}

	if err := b.bundleCSS(assetsDir); err != nil {
		return fmt.Errorf("CSS bundling failed: %w", err)
	}

	if err := b.bundleJS(assetsDir); err != nil {
		return fmt.Errorf("JS bundling failed: %w", err)
	}

	return nil
}

func (b *Bundler) bundleCSS(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	cssFiles := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".css" {
			cssFiles = append(cssFiles, filepath.Join(dir, entry.Name()))
		}
	}

	if len(cssFiles) == 0 {
		fmt.Println("No CSS files to bundle")
		return nil
	}

	fmt.Printf("Bundling %d CSS files...\n", len(cssFiles))

	var wg sync.WaitGroup
	sem := make(chan struct{}, b.parallelism)
	processed := 0

	for _, cssPath := range cssFiles {
		wg.Add(1)
		sem <- struct{}{}

		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := b.processCSS(path); err != nil {
				fmt.Printf("Failed to process %s: %v\n", filepath.Base(path), err)
				return
			}
			processed++
		}(cssPath)
	}

	wg.Wait()
	fmt.Printf("CSS processing complete: %d files processed\n", processed)
	return nil
}

func (b *Bundler) processCSS(cssPath string) error {
	content, err := os.ReadFile(cssPath)
	if err != nil {
		return err
	}

	css := string(content)

	if b.minify {
		css = b.minifyCSS(css)
	}

	outputPath := strings.TrimSuffix(cssPath, filepath.Ext(cssPath)) + ".min.css"
	if err := os.WriteFile(outputPath, []byte(css), 0644); err != nil {
		return err
	}

	originalSize := len(content)
	optimizedSize := len(css)
	savings := float64(originalSize-optimizedSize) / float64(originalSize) * 100

	fmt.Printf("  %s: %d -> %d bytes (%.1f%% reduction)\n",
		filepath.Base(cssPath), originalSize, optimizedSize, savings)

	return nil
}

func (b *Bundler) minifyCSS(css string) string {
	css = regexp.MustCompile(`/\*[\s\S]*?\*/`).ReplaceAllString(css, "")
	css = regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")
	css = regexp.MustCompile(`\s*([{};,:])\s*`).ReplaceAllString(css, "$1")
	css = regexp.MustCompile(`;}`).ReplaceAllString(css, "}")
	css = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(css, "")

	return css
}

func (b *Bundler) bundleJS(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	jsFiles := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".js" {
			jsFiles = append(jsFiles, filepath.Join(dir, entry.Name()))
		}
	}

	if len(jsFiles) == 0 {
		fmt.Println("No JS files to bundle")
		return nil
	}

	fmt.Printf("Bundling %d JS files...\n", len(jsFiles))

	var wg sync.WaitGroup
	sem := make(chan struct{}, b.parallelism)
	processed := 0

	for _, jsPath := range jsFiles {
		wg.Add(1)
		sem <- struct{}{}

		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := b.processJS(path); err != nil {
				fmt.Printf("Failed to process %s: %v\n", filepath.Base(path), err)
				return
			}
			processed++
		}(jsPath)
	}

	wg.Wait()
	fmt.Printf("JS processing complete: %d files processed\n", processed)
	return nil
}

func (b *Bundler) processJS(jsPath string) error {
	content, err := os.ReadFile(jsPath)
	if err != nil {
		return err
	}

	js := string(content)

	if b.minify {
		js = b.minifyJS(js)
	}

	outputPath := strings.TrimSuffix(jsPath, filepath.Ext(jsPath)) + ".min.js"
	if err := os.WriteFile(outputPath, []byte(js), 0644); err != nil {
		return err
	}

	originalSize := len(content)
	optimizedSize := len(js)
	savings := float64(originalSize-optimizedSize) / float64(originalSize) * 100

	fmt.Printf("  %s: %d -> %d bytes (%.1f%% reduction)\n",
		filepath.Base(jsPath), originalSize, optimizedSize, savings)

	return nil
}

func (b *Bundler) minifyJS(js string) string {
	js = regexp.MustCompile(`/\*[\s\S]*?\*/`).ReplaceAllString(js, "")
	js = regexp.MustCompile(`//[^\n]*`).ReplaceAllString(js, "")
	js = regexp.MustCompile(`\s+`).ReplaceAllString(js, " ")
	js = regexp.MustCompile(`\s*([{};,:])\s*`).ReplaceAllString(js, "$1")
	js = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(js, "")

	return js
}

func CombineCSSFiles(files []string, outputPath string) error {
	var buf bytes.Buffer

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}
		buf.Write(content)
		buf.WriteString("\n")
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func CombineJSFiles(files []string, outputPath string) error {
	var buf bytes.Buffer

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}
		buf.Write(content)
		buf.WriteString(";\n")
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}
