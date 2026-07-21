package naabu

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
	"github.com/kartik-Gehlot/RedRecon/internal/worker"
)

func Run(scan *models.Scan) error {

	worker.Run(len(scan.Hosts), 10, func(i int) {

		cmd := exec.Command(
			"naabu",
			"-silent",
			"-host",
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

			line := scanner.Text()

			parts := strings.Split(line, ":")

			if len(parts) != 2 {
				continue
			}

			port, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			scan.Hosts[i].Ports = append(scan.Hosts[i].Ports, port)
		}

		if err := scanner.Err(); err != nil {
			return
		}

		cmd.Wait()
	})
	return nil
}
