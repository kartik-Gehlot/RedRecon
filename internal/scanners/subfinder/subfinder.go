package subfinder

import (
	"bufio"
	"os/exec"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

func Run(scan *models.Scan) error {

	cmd := exec.Command(
		"subfinder",
		"-silent",
		"-d",
		scan.Target,
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

		host := models.Host{
			Hostname: scanner.Text(),
		}

		scan.Hosts = append(scan.Hosts, host)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return cmd.Wait()
}
