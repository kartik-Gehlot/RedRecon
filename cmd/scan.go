package cmd

import (
	"fmt"

	"github.com/kartik-Gehlot/RedRecon/internal/validator"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a target",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Args: %#v\n", args)
		fmt.Printf("Length: %d\n", len(args))
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

	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
