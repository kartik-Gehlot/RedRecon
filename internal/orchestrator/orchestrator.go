package orchestrator

import (
	"fmt"
	"time"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
	subfinder "github.com/kartik-Gehlot/RedRecon/internal/scanners/subfinder"
	"github.com/kartik-Gehlot/RedRecon/internal/validator"
)

func Start(target string) error {

	scan := models.Scan{
		Target:    target,
		Status:    "Running",
		Version:   "v0.1.0",
		StartTime: time.Now(),
	}

	fmt.Println("[INFO] Starting RedRecon")

	// Validate Target
	if err := validator.ValidateTarget(scan.Target); err != nil {
		return err
	}

	fmt.Println("[INFO] Running Subfinder...")

	// Run Subfinder
	if err := subfinder.Run(&scan); err != nil {
		return err
	}

	// Scan Completed
	scan.EndTime = time.Now()
	scan.Status = "Completed"

	fmt.Printf("\nFound %d hosts\n\n", len(scan.Hosts))

	for _, host := range scan.Hosts {
		fmt.Println(host.Hostname)
	}

	return nil
}
