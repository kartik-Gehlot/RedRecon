package httpx

import (
	"bufio"
	"encoding/json"
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

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		if err := cmd.Start(); err != nil {
			return err
		}

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {

			var result Result

			if err := json.Unmarshal(scanner.Bytes(), &result); err != nil {
				continue
			}

			if result.Failed {
				continue
			}

			scan.Hosts[i].Alive = true
			scan.Hosts[i].IP = result.HostIP
			scan.Hosts[i].StatusCode = result.StatusCode
			scan.Hosts[i].Title = result.Title
			scan.Hosts[i].Technologies = result.Tech
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		cmd.Wait()
	}

	return nil
}
