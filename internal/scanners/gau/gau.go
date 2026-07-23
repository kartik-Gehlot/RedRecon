package gau

import (
	"bufio"
	"fmt"

	"os/exec"

	"github.com/kartik-Gehlot/RedRecon/internal/config"
	"github.com/kartik-Gehlot/RedRecon/internal/logger"
	"github.com/kartik-Gehlot/RedRecon/internal/models"
	"github.com/kartik-Gehlot/RedRecon/internal/worker"
)

func Run(scan *models.Scan) error {

	worker.Run(len(scan.Hosts), config.Workers, func(i int) {

		if !scan.Hosts[i].Alive {
			return
		}

		cmd := exec.Command(
			"gau",
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
		count := 0

		for scanner.Scan() {

			url := scanner.Text()

			scan.Hosts[i].URLs = append(scan.Hosts[i].URLs, url)

			count++

			if count%100 == 0 {
				logger.Info(fmt.Sprintf("%s : %d URLs collected", scan.Hosts[i].Hostname, count))
			}
		}

		logger.Success(fmt.Sprintf("%s : %d URLs collected", scan.Hosts[i].Hostname, count))

		if err := scanner.Err(); err != nil {
			return
		}

		cmd.Wait()

	})

	return nil
}
