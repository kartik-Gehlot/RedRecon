package katana

import (
	"bufio"
	"os/exec"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

func Run(scan *models.Scan) error {

	for i := range scan.Hosts {

		if !scan.Hosts[i].Alive {
			continue
		}

		cmd := exec.Command(
			"katana",
			"-silent",
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
			scan.Hosts[i].URLs = append(
				scan.Hosts[i].URLs,
				scanner.Text(),
			)
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		cmd.Wait()
	}

	return nil
}
