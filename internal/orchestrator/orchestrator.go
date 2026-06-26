package orchestrator

import (
	"fmt"
	"time"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
	httpx "github.com/kartik-Gehlot/RedRecon/internal/scanners/httpx"
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

	if err := validator.ValidateTarget(scan.Target); err != nil {
		return err
	}

	fmt.Println("[INFO] Running Subfinder...")

	if err := subfinder.Run(&scan); err != nil {
		return err
	}
	fmt.Println("[INFO] Running Httpx...")

	if err := httpx.Run(&scan); err != nil {
		return err
	}
	scan.EndTime = time.Now()
	scan.Status = "Completed"

	fmt.Printf("\nFound %d hosts\n\n", len(scan.Hosts))

	for _, host := range scan.Hosts {

		fmt.Println("--------------------------------")
		fmt.Println("Host        :", host.Hostname)
		fmt.Println("Alive       :", host.Alive)
		fmt.Println("IP          :", host.IP)
		fmt.Println("Status Code :", host.StatusCode)
		fmt.Println("Title       :", host.Title)

		if len(host.Technologies) > 0 {
			fmt.Println("Technologies:", host.Technologies)
		}

		fmt.Println("--------------------------------")
	}

	return nil
}
