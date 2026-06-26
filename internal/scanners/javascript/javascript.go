package javascript

import (
	"strings"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

func Run(scan *models.Scan) error {

	for i := range scan.Hosts {

		for _, url := range scan.Hosts[i].URLs {

			if strings.HasSuffix(url, ".js") {
				scan.Hosts[i].JavaScripts = append(
					scan.Hosts[i].JavaScripts,
					url,
				)
			}
		}
	}

	return nil
}
