package orchestrator

import (
	"fmt"
	"time"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
	gau "github.com/kartik-Gehlot/RedRecon/internal/scanners/gau"
	httpx "github.com/kartik-Gehlot/RedRecon/internal/scanners/httpx"
	javascript "github.com/kartik-Gehlot/RedRecon/internal/scanners/javascript"
	katana "github.com/kartik-Gehlot/RedRecon/internal/scanners/katana"
	naabu "github.com/kartik-Gehlot/RedRecon/internal/scanners/naabu"
	"github.com/kartik-Gehlot/RedRecon/internal/scanners/nuclei"
	secrets "github.com/kartik-Gehlot/RedRecon/internal/scanners/secrets"
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
	fmt.Println("[INFO] Running Naabu...")

	if err := naabu.Run(&scan); err != nil {
		return err
	}
	fmt.Println("[INFO] Running Katana...")

	if err := katana.Run(&scan); err != nil {
		return err
	}
	fmt.Println("[INFO] Running GAU...")

	if err := gau.Run(&scan); err != nil {
		return err
	}
	fmt.Println("[INFO] Collecting JavaScript files...")

	if err := javascript.Run(&scan); err != nil {
		return err
	}
	fmt.Println("[INFO] Scanning JavaScript for Secrets...")

	if err := secrets.Run(&scan); err != nil {
		return err
	}
	fmt.Println("[INFO] Running Nuclei...")

	if err := nuclei.Run(&scan); err != nil {
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
		fmt.Println("Ports       :", host.Ports)
		fmt.Println("URLs Found   :", len(host.URLs))
		fmt.Println("JavaScript Files :", len(host.JavaScripts))
		fmt.Println("Secrets Found :", len(host.Secrets))
		fmt.Println("Findings :", len(host.Findings))

		if len(host.Technologies) > 0 {
			fmt.Println("Technologies:", host.Technologies)
		}

		fmt.Println("--------------------------------")
	}

	return nil
}
