package dashboard

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

func Start() error {

	file, err := os.Open("output/scan.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var scan models.Scan

	if err := json.NewDecoder(file).Decode(&scan); err != nil {
		return err
	}

	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Build(&scan)

		tmpl.ExecuteTemplate(w, "dashboard", data)
	})

	return http.ListenAndServe(":8080", nil)
}
