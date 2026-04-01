package cmd

import (
	"fmt"
	"os"

	"github.com/pangolin-cms/staticpress/cmd/internal/config"
	"github.com/pangolin-cms/staticpress/cmd/internal/exporter"

	"github.com/spf13/cobra"
)

var (
	deployPlatform string
	netlifyToken   string
	netlifySite    string
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy static files to S3 or Netlify",
	Long:  `Upload the exported static files to S3 or Netlify.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		distDir, _ := cmd.Flags().GetString("dist")

		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		switch deployPlatform {
		case "s3":
			return deployToS3(cmd, cfg, distDir)
		case "netlify":
			return deployToNetlify(cmd, cfg, distDir)
		default:
			return fmt.Errorf("invalid platform: %s (use 's3' or 'netlify')", deployPlatform)
		}
	},
}

func deployToS3(cmd *cobra.Command, cfg *config.Config, distDir string) error {
	bucket, _ := cmd.Flags().GetString("bucket")
	region, _ := cmd.Flags().GetString("region")

	if bucket == "" {
		return fmt.Errorf("please provide --bucket flag for S3 deployment")
	}

	if err := exporter.DeployToS3(distDir, bucket, region, cfg); err != nil {
		return fmt.Errorf("deploy failed: %w", err)
	}

	fmt.Printf("Successfully deployed to s3://%s\n", bucket)
	return nil
}

func deployToNetlify(cmd *cobra.Command, cfg *config.Config, distDir string) error {
	token := netlifyToken
	if token == "" {
		token = os.Getenv("NETLIFY_AUTH_TOKEN")
	}
	if token == "" {
		token = cfg.NetlifyToken
	}
	if token == "" {
		return fmt.Errorf("Netlify auth token required. Set --netlify-token flag, NETLIFY_AUTH_TOKEN env var, or 'netlify_token' in config")
	}

	site := netlifySite
	if site == "" {
		site = cfg.NetlifySite
	}
	if site == "" {
		return fmt.Errorf("Netlify site ID required. Set --netlify-site flag or 'netlify_site' in config")
	}

	deployer := exporter.NewNetlifyDeployer(token, site, distDir)
	if err := deployer.Deploy(); err != nil {
		return fmt.Errorf("deploy failed: %w", err)
	}

	return nil
}

func init() {
	DeployCmd.Flags().StringP("dist", "d", "dist", "Directory to deploy")
	DeployCmd.Flags().StringP("platform", "p", "s3", "Deployment platform: s3 or netlify")
	DeployCmd.Flags().StringP("bucket", "b", "", "S3 bucket name (required for s3 platform)")
	DeployCmd.Flags().StringP("region", "r", "us-east-1", "AWS region")
	DeployCmd.Flags().StringVar(&netlifyToken, "netlify-token", "", "Netlify auth token")
	DeployCmd.Flags().StringVar(&netlifySite, "netlify-site", "", "Netlify site ID or name")
}
