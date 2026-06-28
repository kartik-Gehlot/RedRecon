package output

import (
	"encoding/json"
	"os"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

func Save(scan *models.Scan) error {

	file, err := os.Create("output/scan.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(scan)
}
