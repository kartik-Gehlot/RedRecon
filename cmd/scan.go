package cmd

import (
	"fmt"

	"github.com/kartik-Gehlot/RedRecon/internal/orchestrator"
	"github.com/kartik-Gehlot/RedRecon/internal/validator"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a target",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("Please provide a target")
			return
		}

		target := args[0]

		err := validator.ValidateTarget(target)

		if err != nil {
			fmt.Println("❌", err)
			return
		}

		fmt.Println("✅ Target Valid:", target)
		err = orchestrator.Start(target)
		if err != nil {
			fmt.Println("❌", err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
