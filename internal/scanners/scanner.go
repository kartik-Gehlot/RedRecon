package scanners

import "github.com/kartik-Gehlot/RedRecon/internal/models"

type Scanner interface {
	Run(scan *models.Scan) error
}
