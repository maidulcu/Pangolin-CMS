package exporter

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

type ImageOptimizer struct {
	quality     int
	format      string
	enabled     bool
	parallelism int
}

type OptimizeOptions struct {
	Quality     int
	Format      string
	Enabled     bool
	Parallelism int
}

func NewImageOptimizer(opts OptimizeOptions) *ImageOptimizer {
	quality := opts.Quality
	if quality <= 0 {
		quality = 80
	}

	format := strings.ToLower(opts.Format)
	if format == "" {
		format = "webp"
	}

	parallelism := opts.Parallelism
	if parallelism <= 0 {
		parallelism = 4
	}

	return &ImageOptimizer{
		quality:     quality,
		format:      format,
		enabled:     opts.Enabled,
		parallelism: parallelism,
	}
}

func (o *ImageOptimizer) OptimizeDirectory(dir string) error {
	if !o.enabled {
		return nil
	}

	imagesDir := filepath.Join(dir, "images")
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		fmt.Println("No images directory found, skipping optimization")
		return nil
	}

	entries, err := os.ReadDir(imagesDir)
	if err != nil {
		return err
	}

	imageFiles := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" {
			imageFiles = append(imageFiles, filepath.Join(imagesDir, entry.Name()))
		}
	}

	if len(imageFiles) == 0 {
		fmt.Println("No images to optimize")
		return nil
	}

	fmt.Printf("Optimizing %d images to %s (quality: %d)...\n", len(imageFiles), o.format, o.quality)

	var wg sync.WaitGroup
	sem := make(chan struct{}, o.parallelism)
	successCount := 0
	errorCount := 0

	for _, imgPath := range imageFiles {
		wg.Add(1)
		sem <- struct{}{}

		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := o.optimizeImage(path); err != nil {
				fmt.Printf("Failed to optimize %s: %v\n", filepath.Base(path), err)
				errorCount++
				return
			}
			successCount++
		}(imgPath)
	}

	wg.Wait()

	fmt.Printf("Image optimization complete: %d optimized, %d failed\n", successCount, errorCount)
	return nil
}

func (o *ImageOptimizer) optimizeImage(imgPath string) error {
	file, err := os.Open(imgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	srcImage, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	resized := imaging.Resize(srcImage, 0, 0, imaging.Lanczos)

	var outputPath string
	switch o.format {
	case "webp":
		outputPath = strings.TrimSuffix(imgPath, filepath.Ext(imgPath)) + ".webp"
	case "avif":
		outputPath = strings.TrimSuffix(imgPath, filepath.Ext(imgPath)) + ".avif"
	default:
		return fmt.Errorf("unsupported format: %s", o.format)
	}

	switch o.format {
	case "webp":
		if err := o.encodeWebP(resized, outputPath); err != nil {
			return err
		}
	case "avif":
		if err := o.encodeAVIF(resized, outputPath); err != nil {
			return err
		}
	}

	originalSize := getFileSize(imgPath)
	optimizedSize := getFileSize(outputPath)
	savings := float64(originalSize-optimizedSize) / float64(originalSize) * 100

	fmt.Printf("  %s: %s -> %s (%.1f%% reduction)\n",
		filepath.Base(imgPath), formatFileSize(originalSize), formatFileSize(optimizedSize), savings)

	if imgPath != outputPath {
		os.Remove(imgPath)
	}

	return nil
}

func (o *ImageOptimizer) encodeWebP(img image.Image, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	opts := &webp.Options{
		Quality: float32(o.quality),
	}

	return webp.Encode(file, img, opts)
}

func (o *ImageOptimizer) encodeAVIF(img image.Image, outputPath string) error {
	return fmt.Errorf("AVIF encoding not yet supported, please use WebP format")
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func DetectBestFormat(imgPath string) (string, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return "", err
	}

	return format, nil
}

func CalculateSavings(originalPath, optimizedPath string) (float64, error) {
	origSize := getFileSize(originalPath)
	optSize := getFileSize(optimizedPath)

	if origSize == 0 {
		return 0, fmt.Errorf("original file size is zero")
	}

	return float64(origSize-optSize) / float64(origSize) * 100, nil
}
