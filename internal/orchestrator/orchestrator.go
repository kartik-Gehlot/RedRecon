package orchestrator

import (
	"fmt"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
	"github.com/kartik-Gehlot/RedRecon/internal/modules/subfinder"
)

func Start(target string) error {
	fmt.Println("\n🚀 Starting RedRecon Scan")

	// Create Scan object
	scan := &models.Scan{
		Target: target,
	}

	// Run Subfinder
	subdomains, err := subfinder.Run(scan.Target)
	if err != nil {
		return err
	}

	// Save results
	scan.Subdomains = subdomains

	// Print results
	fmt.Printf("\n✅ Found %d subdomains\n\n", len(scan.Subdomains))

	for _, sub := range scan.Subdomains {
		fmt.Println(sub)
	}

	return nil
}
