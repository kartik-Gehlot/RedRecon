package summary

import (
	"fmt"
	"time"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

func Print(scan *models.Scan) {

	var alive int
	var ports int
	var urls int
	var findings int

	for _, host := range scan.Hosts {

		if host.Alive {
			alive++
		}

		ports += len(host.Ports)
		urls += len(host.URLs)
		findings += len(host.Findings)
	}

	fmt.Println()
	fmt.Println("==============================================")
	fmt.Println("          REDRECON SCAN SUMMARY")
	fmt.Println("==============================================")

	fmt.Printf("Target          : %s\n", scan.Target)
	fmt.Printf("Version         : %s\n", scan.Version)
	fmt.Printf("Status          : %s\n", scan.Status)
	fmt.Printf("Duration        : %v\n", scan.EndTime.Sub(scan.StartTime).Round(time.Second))

	fmt.Println("----------------------------------------------")

	fmt.Printf("Hosts Found     : %d\n", len(scan.Hosts))
	fmt.Printf("Alive Hosts     : %d\n", alive)
	fmt.Printf("Open Ports      : %d\n", ports)
	fmt.Printf("URLs Collected  : %d\n", urls)
	fmt.Printf("Findings        : %d\n", findings)
	fmt.Printf("GitHub Findings : %d\n", len(scan.GithubFindings))

	fmt.Println("----------------------------------------------")

	fmt.Println("JSON Report : output/scan.json")
	fmt.Println("HTML Report : output/report.html")
	fmt.Println("PDF Report  : output/report.pdf")

	fmt.Println("==============================================")
}
