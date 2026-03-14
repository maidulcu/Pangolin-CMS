package cmd

import (
	"fmt"

	"github.com/pangolin-cms/staticpress/cmd/internal/config"
	"github.com/pangolin-cms/staticpress/cmd/internal/exporter"

	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy static files to S3",
	Long:  `Upload the exported static files to S3 and optionally invalidate CloudFront cache.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		distDir, _ := cmd.Flags().GetString("dist")
		bucket, _ := cmd.Flags().GetString("bucket")
		region, _ := cmd.Flags().GetString("region")

		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if bucket == "" {
			return fmt.Errorf("please provide --bucket flag")
		}

		if err := exporter.DeployToS3(distDir, bucket, region, cfg); err != nil {
			return fmt.Errorf("deploy failed: %w", err)
		}

		fmt.Printf("Successfully deployed to s3://%s\n", bucket)
		return nil
	},
}

func init() {
	DeployCmd.Flags().StringP("dist", "d", "dist", "Directory to deploy")
	DeployCmd.Flags().StringP("bucket", "b", "", "S3 bucket name")
	DeployCmd.Flags().StringP("region", "r", "us-east-1", "AWS region")
}
