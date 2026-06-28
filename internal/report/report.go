package report

import (
	"html/template"
	"os"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

func Generate(scan *models.Scan) error {

	tmpl, err := template.ParseFiles("templates/report.html")
	if err != nil {
		return err
	}

	file, err := os.Create("output/report.html")
	if err != nil {
		return err
	}
	defer file.Close()

	data := Build(scan)

	return tmpl.Execute(file, data)
}
