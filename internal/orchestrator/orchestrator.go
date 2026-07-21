package orchestrator

import (
	"fmt"
	"sync"
	"time"

	"github.com/kartik-Gehlot/RedRecon/internal/ai"
	"github.com/kartik-Gehlot/RedRecon/internal/logger"
	"github.com/kartik-Gehlot/RedRecon/internal/models"
	"github.com/kartik-Gehlot/RedRecon/internal/output"
	"github.com/kartik-Gehlot/RedRecon/internal/pdf"
	"github.com/kartik-Gehlot/RedRecon/internal/report"
	gau "github.com/kartik-Gehlot/RedRecon/internal/scanners/gau"
	"github.com/kartik-Gehlot/RedRecon/internal/stats"
	"github.com/kartik-Gehlot/RedRecon/internal/summary"

	github "github.com/kartik-Gehlot/RedRecon/internal/scanners/github"
	httpx "github.com/kartik-Gehlot/RedRecon/internal/scanners/httpx"
	javascript "github.com/kartik-Gehlot/RedRecon/internal/scanners/javascript"
	katana "github.com/kartik-Gehlot/RedRecon/internal/scanners/katana"
	naabu "github.com/kartik-Gehlot/RedRecon/internal/scanners/naabu"
	"github.com/kartik-Gehlot/RedRecon/internal/scanners/nuclei"

	//	secrets "github.com/kartik-Gehlot/RedRecon/internal/scanners/secrets"
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

	logger.Info("Starting RedRecon")

	if err := validator.ValidateTarget(scan.Target); err != nil {
		return err
	}
	start := time.Now()
	logger.Info("Running Subfinder")

	if err := subfinder.Run(&scan); err != nil {
		return err
	}
	stats.Record("Subfinder", start)
	logger.Success("Subfinder completed")
	start = time.Now()
	logger.Info("Running Httpx & Naabu")

	var wg sync.WaitGroup

	wg.Add(2)

	var httpxErr error
	var naabuErr error

	go func() {
		defer wg.Done()
		httpxErr = httpx.Run(&scan)
	}()

	go func() {
		defer wg.Done()
		naabuErr = naabu.Run(&scan)
	}()

	wg.Wait()

	if httpxErr != nil {
		return httpxErr
	}

	if naabuErr != nil {
		return naabuErr
	}
	logger.Success("Httpx & Naabu completed")
	stats.Record("Httpx & Naabu", start)

	start = time.Now()
	logger.Info("Running Katana & GAU...")

	wg.Add(2)

	var katanaErr error
	var gauErr error

	go func() {
		defer wg.Done()
		katanaErr = katana.Run(&scan)
	}()

	go func() {
		defer wg.Done()
		gauErr = gau.Run(&scan)
	}()

	wg.Wait()

	if katanaErr != nil {
		return katanaErr
	}

	if gauErr != nil {
		return gauErr
	}
	logger.Success("Katana & GAU completed")
	stats.Record("Katana & GAU", start)
	start = time.Now()
	logger.Info("Running JavaScript & GitHub Search...")

	wg.Add(2)

	var jsErr error
	var githubErr error

	go func() {
		defer wg.Done()
		jsErr = javascript.Run(&scan)
	}()

	go func() {
		defer wg.Done()
		githubErr = github.Run(&scan)
	}()

	wg.Wait()

	if jsErr != nil {
		return jsErr
	}

	if githubErr != nil {
		return githubErr
	}
	logger.Success("JavaScript & GitHub Search completed")
	stats.Record("JavaScript & GitHub Search", start)

	start = time.Now()
	logger.Info("Running Nuclei...")

	if err := nuclei.Run(&scan); err != nil {
		return err
	}
	stats.Record("Nuclei", start)
	logger.Success("Nuclei completed")
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
		// fmt.Println("Secrets Found :", len(host.Secrets))
		fmt.Println("Findings :", len(host.Findings))

		if len(host.Technologies) > 0 {
			fmt.Println("Technologies:", host.Technologies)
		}

		fmt.Println("--------------------------------")
	}
	fmt.Printf("\nGitHub Findings: %d\n", len(scan.GithubFindings))

	for _, finding := range scan.GithubFindings {
		fmt.Println("--------------------------------")
		fmt.Println("Repository :", finding.Repository)
		fmt.Println("File       :", finding.File)
		fmt.Println("URL        :", finding.URL)
	}
	start = time.Now()
	logger.Info("Saving Scan Results")

	if err := output.Save(&scan); err != nil {
		return err
	}
	stats.Record("Save Results", start)
	logger.Success("Scan results saved")
	start = time.Now()
	logger.Info("Generating HTML Report...")

	if err := report.Generate(&scan); err != nil {
		return err
	}
	stats.Record("Generate HTML Report", start)
	logger.Success("HTML Report generated")
	start = time.Now()
	logger.Info("Running AI Analysis...")

	analysis := ai.Analyze(&scan)
	stats.Record("AI Analysis", start)
	logger.Success("AI Analysis completed")
	start = time.Now()
	logger.Info("Generating PDF Report...")

	if err := pdf.Generate(&scan, analysis); err != nil {
		return err
	}
	stats.Record("Generate PDF Report", start)
	logger.Success("PDF Report generated")

	stats.Print()
	summary.Print(&scan)
	return nil
}
