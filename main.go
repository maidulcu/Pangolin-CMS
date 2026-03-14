package main

import (
	"fmt"
	"os"

	"staticpress/cmd"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "staticpress",
		Short: "StaticPress - Export WordPress sites to static HTML",
		Long:  `A CLI tool to export WordPress sites to static HTML files for deployment to S3, Netlify, or other static hosting providers.`,
	}

	rootCmd.AddCommand(cmd.InitCmd)
	rootCmd.AddCommand(cmd.ExportCmd)
	rootCmd.AddCommand(cmd.DeployCmd)
	rootCmd.AddCommand(cmd.ServeCmd)
	rootCmd.AddCommand(cmd.DashboardCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
