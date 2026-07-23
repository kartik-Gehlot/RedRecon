package httpx

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
	"github.com/kartik-Gehlot/RedRecon/internal/worker"
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

	worker.Run(len(scan.Hosts), 10, func(i int) {

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
			fmt.Println("httpx error:", err)
			return
		}

		if len(output) == 0 {
			return
		}

		var result Result
		if err := json.Unmarshal(output, &result); err != nil {
			fmt.Println("JSON error:", err)
			return
		}

		if result.Failed {
			return
		}

		scan.Hosts[i].Alive = true
		scan.Hosts[i].IP = result.HostIP
		scan.Hosts[i].StatusCode = result.StatusCode
		scan.Hosts[i].Title = result.Title
		scan.Hosts[i].Technologies = result.Tech
	})

	return nil
}
