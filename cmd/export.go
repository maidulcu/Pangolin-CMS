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

		exp := exporter.NewExporter(distDir, concurrency)
		if err := exp.Export(urls); err != nil {
			return fmt.Errorf("export failed: %w", err)
		}

		fmt.Printf("Successfully exported to %s\n", distDir)

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

		return nil
	},
}

func init() {
	ExportCmd.Flags().IntP("concurrency", "c", 5, "Number of concurrent requests")
	ExportCmd.Flags().StringP("dist", "d", "dist", "Output directory")
	ExportCmd.Flags().BoolVar(&optImages, "optimize-images", false, "Enable image optimization")
	ExportCmd.Flags().StringVar(&optImageFormat, "image-format", "webp", "Image output format (webp, avif)")
	ExportCmd.Flags().IntVar(&optImageQuality, "image-quality", 80, "Image quality (1-100)")
}
