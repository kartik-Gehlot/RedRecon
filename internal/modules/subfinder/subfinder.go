package subfinder

import (
	"bufio"
	"os/exec"
)

func Run(domain string) ([]string, error) {

	var results []string

	cmd := exec.Command(
		"subfinder",
		"-silent",
		"-d",
		domain,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}
