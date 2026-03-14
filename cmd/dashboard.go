package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var DashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Start the web dashboard",
	Long:  `Start the Pangolin web dashboard for managing exports.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dashboardBinary := "pangolin-dashboard"

		if _, err := os.Stat(dashboardBinary); err != nil {
			dashboardBinary = "./dashboard"
		}

		execObj := exec.Command(dashboardBinary)
		execObj.Stdout = os.Stdout
		execObj.Stderr = os.Stderr
		execObj.Stdin = os.Stdin

		return execObj.Run()
	},
}
