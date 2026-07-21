package katana

import (
	"bufio"
	"os/exec"

	"github.com/kartik-Gehlot/RedRecon/internal/config"
	"github.com/kartik-Gehlot/RedRecon/internal/models"
	"github.com/kartik-Gehlot/RedRecon/internal/worker"
)

func Run(scan *models.Scan) error {
	worker.Run(len(scan.Hosts), config.Workers, func(i int) {

		cmd := exec.Command(
			"katana",
			"-silent",
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
			scan.Hosts[i].URLs = append(
				scan.Hosts[i].URLs,
				scanner.Text(),
			)
		}

		if err := scanner.Err(); err != nil {
			return
		}

		cmd.Wait()
	})

	return nil
}
