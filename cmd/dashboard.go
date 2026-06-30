package cmd

import (
	"fmt"

	"github.com/kartik-Gehlot/RedRecon/internal/dashboard"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{

	Use:   "dashboard",
	Short: "Launch the RedRecon Dashboard",

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("[INFO] Starting Dashboard...")
		fmt.Println("[INFO] Open http://localhost:8080")

		if err := dashboard.Start(); err != nil {
			fmt.Println("[ERROR]", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
