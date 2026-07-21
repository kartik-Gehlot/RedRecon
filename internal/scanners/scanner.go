package scanners

import "github.com/kartik-Gehlot/RedRecon/internal/models"

type Scanner interface {
	Name() string
	Run(*models.Scan) error
}
