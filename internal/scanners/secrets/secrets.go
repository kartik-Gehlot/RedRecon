package secrets

import (
	"io"
	"net/http"
	"regexp"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

var patterns = []*regexp.Regexp{
	regexp.MustCompile(`AIza[0-9A-Za-z\-_]{35}`),   // Google API Key
	regexp.MustCompile(`AKIA[0-9A-Z]{16}`),         // AWS Access Key
	regexp.MustCompile(`ghp_[A-Za-z0-9]{36}`),      // GitHub Token
	regexp.MustCompile(`xox[baprs]-[A-Za-z0-9-]+`), // Slack Token
}

func Run(scan *models.Scan) error {

	client := &http.Client{}

	for i := range scan.Hosts {

		for _, js := range scan.Hosts[i].JavaScripts {

			resp, err := client.Get(js)
			if err != nil {
				continue
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			content := string(body)

			for _, pattern := range patterns {

				matches := pattern.FindAllString(content, -1)

				scan.Hosts[i].Secrets = append(
					scan.Hosts[i].Secrets,
					matches...,
				)
			}
		}
	}

	return nil
}
