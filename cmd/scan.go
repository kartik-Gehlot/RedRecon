package cmd

import (
	"fmt"

	"github.com/kartik-Gehlot/RedRecon/internal/orchestrator"
	"github.com/spf13/cobra"
)

// var profile string
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a target",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a target")
			return
		}

		if err := orchestrator.Start(args[0]); err != nil {
			fmt.Println("[ERROR]", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	
}
