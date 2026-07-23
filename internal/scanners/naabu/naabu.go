package httpx

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

type Result struct {
	Host       string   `json:"host"`
	HostIP     string   `json:"host_ip"`
	StatusCode int      `json:"status_code"`
	Title      string   `json:"title"`
	Tech       []string `json:"tech"`
	Failed     bool     `json:"failed"`
}

func Run(scan *models.Scan) error {

	for i := range scan.Hosts {

		cmd := exec.Command(
			"httpx",
			"-silent",
			"-json",
			"-title",
			"-tech-detect",
			"-u",
			scan.Hosts[i].Hostname,
		)

		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("httpx failed for %s: %v\n", scan.Hosts[i].Hostname, err)
			continue
		}

		if len(output) == 0 {
			fmt.Printf("No output for %s\n", scan.Hosts[i].Hostname)
			continue
		}

		var result Result
		if err := json.Unmarshal(output, &result); err != nil {
			fmt.Printf("JSON error for %s: %v\n", scan.Hosts[i].Hostname, err)
			continue
		}

		fmt.Printf("Alive: %s (%d)\n", scan.Hosts[i].Hostname, result.StatusCode)

		scan.Hosts[i].Alive = true
		scan.Hosts[i].IP = result.HostIP
		scan.Hosts[i].StatusCode = result.StatusCode
		scan.Hosts[i].Title = result.Title
		scan.Hosts[i].Technologies = result.Tech
	}

	return nil
}
