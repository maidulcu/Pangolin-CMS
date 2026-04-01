package cmd

import (
	"fmt"

	"github.com/pangolin-cms/staticpress/cmd/internal/exporter"
	"github.com/pangolin-cms/staticpress/cmd/internal/sitemap"

	"github.com/spf13/cobra"
)

var (
	optImages       bool
	optImageFormat  string
	optImageQuality int
	optMinify       bool
	optIncremental  bool
	optClearCache   bool
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export WordPress site to static HTML",
	Long:  `Crawl your WordPress site and export all pages to static HTML files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		distDir, _ := cmd.Flags().GetString("dist")

		urls, err := sitemap.FetchSitemaps()
		if err != nil {
			return fmt.Errorf("failed to fetch sitemaps: %w", err)
		}

		if len(urls) == 0 {
			fmt.Println("No URLs found in sitemap")
			return nil
		}

		fmt.Printf("Found %d URLs to export\n", len(urls))

		if optClearCache {
			fmt.Println("Clearing export cache...")
			incExp := exporter.NewIncrementalExporter(distDir, concurrency)
			if err := incExp.ClearCache(); err != nil {
				fmt.Printf("Cache cleared: %v\n", err)
			} else {
				fmt.Println("Cache cleared successfully")
			}
		}

		var totalExported int
		if optIncremental {
			incExp := exporter.NewIncrementalExporter(distDir, concurrency)
			totalExported, err = incExp.ExportIncremental(urls)
			if err != nil {
				return fmt.Errorf("incremental export failed: %w", err)
			}
		} else {
			exp := exporter.NewExporter(distDir, concurrency)
			if err := exp.Export(urls); err != nil {
				return fmt.Errorf("export failed: %w", err)
			}
			totalExported = len(urls)
		}

		if !optIncremental || totalExported > 0 {
			fmt.Printf("Successfully exported %d pages to %s\n", totalExported, distDir)
		}

		if optImages {
			fmt.Println("\nOptimizing images...")
			optimizer := exporter.NewImageOptimizer(exporter.OptimizeOptions{
				Enabled:     true,
				Format:      optImageFormat,
				Quality:     optImageQuality,
				Parallelism: concurrency,
			})

			if err := optimizer.OptimizeDirectory(distDir); err != nil {
				return fmt.Errorf("image optimization failed: %w", err)
			}
		}

		if optMinify {
			fmt.Println("\nMinifying assets...")
			bundler := exporter.NewBundler(exporter.BundlerOptions{
				Minify:      true,
				Parallelism: concurrency,
			})

			if err := bundler.BundleDirectory(distDir); err != nil {
				return fmt.Errorf("asset minification failed: %w", err)
			}
		}

		return nil
	},
}

func init() {
	ExportCmd.Flags().IntP("concurrency", "c", 5, "Number of concurrent requests")
	ExportCmd.Flags().StringP("dist", "d", "dist", "Output directory")
	ExportCmd.Flags().BoolVar(&optImages, "optimize-images", false, "Enable image optimization")
	ExportCmd.Flags().StringVar(&optImageFormat, "image-format", "webp", "Image output format (webp, avif)")
	ExportCmd.Flags().IntVar(&optImageQuality, "image-quality", 80, "Image quality (1-100)")
	ExportCmd.Flags().BoolVar(&optMinify, "minify", false, "Enable CSS/JS minification")
	ExportCmd.Flags().BoolVar(&optIncremental, "incremental", false, "Only export changed pages (uses ETag/Last-Modified)")
	ExportCmd.Flags().BoolVar(&optClearCache, "clear-cache", false, "Clear the export cache before exporting")
}
