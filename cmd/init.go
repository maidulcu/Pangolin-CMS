package cmd

import (
	"fmt"

	"github.com/pangolin-cms/staticpress/cmd/internal/config"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Pangolin configuration",
	Long:  `Initialize Pangolin by connecting to your WordPress site and saving the configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		apiKey, _ := cmd.Flags().GetString("api-key")

		if url == "" {
			return fmt.Errorf("please provide --url flag or run interactively")
		}

		if apiKey == "" {
			return fmt.Errorf("please provide --api-key flag")
		}

		cfg := &config.Config{
			SiteURL: url,
			APIKey:  apiKey,
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println("Pangolin initialized successfully!")
		fmt.Printf("Site URL: %s\n", url)
		return nil
	},
}

func init() {
	InitCmd.Flags().StringP("url", "u", "", "WordPress site URL")
	InitCmd.Flags().StringP("api-key", "k", "", "API key from WP plugin")
}
