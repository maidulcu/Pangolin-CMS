package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a local server to preview exported site",
	Long:  `Start an HTTP server to preview the exported static site.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		distDir, _ := cmd.Flags().GetString("dist")
		port, _ := cmd.Flags().GetInt("port")

		if _, err := os.Stat(distDir); os.IsNotExist(err) {
			return fmt.Errorf("directory %s does not exist, run 'pangolin export' first", distDir)
		}

		fmt.Printf("Starting server at http://localhost:%d\n", port)
		fmt.Printf("Serving files from: %s\n", distDir)
		fmt.Println("Press Ctrl+C to stop")

		return http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(distDir)))
	},
}

func init() {
	ServeCmd.Flags().StringP("dist", "d", "dist", "Directory to serve")
	ServeCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
}
