package httpx

import (
	"bufio"
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

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return
		}

		if err := cmd.Start(); err != nil {
			return
		}

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {

    fmt.Println("RAW:", scanner.Text())

    var result Result

    if err := json.Unmarshal(scanner.Bytes(), &result); err != nil {
        fmt.Println("JSON ERROR:", err)
        continue
    }

    fmt.Printf("PARSED: %+v\n", result)

    if result.Failed {
        fmt.Println("FAILED RESULT")
        continue
    }

    fmt.Println("SETTING ALIVE:", scan.Hosts[i].Hostname)

    scan.Hosts[i].Alive = true
    scan.Hosts[i].IP = result.HostIP
    scan.Hosts[i].StatusCode = result.StatusCode
    scan.Hosts[i].Title = result.Title
    scan.Hosts[i].Technologies = result.Tech
}

		
		fmt.Println("Scanning:", scan.Hosts[i].Hostname)
		if err := cmd.Wait(); err != nil {
			fmt.Println("Wait error:", err)
		}if err := scanner.Err(); err != nil {

			    fmt.Println("Scanner error:", err)
				
		}
	})
	return nil
}
