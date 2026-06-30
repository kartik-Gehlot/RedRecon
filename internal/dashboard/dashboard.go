package dashboard

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/kartik-Gehlot/RedRecon/internal/ai"
	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

type DashboardPage struct {
	DashboardData
	AI ai.Analysis
}

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

	// data := Build(&scan)

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// dashboardData := Build(&scan)
		// analysis := ai.Analyze(&scan)
		analysis := ai.Analyze(&scan)

		scan.RiskScore = analysis.RiskScore

		dashboardData := Build(&scan)
		page := DashboardPage{
			DashboardData: dashboardData,
			AI:            analysis,
		}

		if err := tmpl.Execute(w, page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	// http.Handle("/report.pdf",
	// 	http.StripPrefix("/",
	// 		http.FileServer(http.Dir("output")),
	// 	),
	// )
	// // Serve output files
	// http.Handle("/scan.json", http.StripPrefix("/", http.FileServer(http.Dir("output"))))
	// http.Handle("/report.html", http.StripPrefix("/", http.FileServer(http.Dir("output"))))

	http.Handle("/scan.json",
		http.StripPrefix("/", http.FileServer(http.Dir("output"))))

	http.Handle("/report.html",
		http.StripPrefix("/", http.FileServer(http.Dir("output"))))

	http.Handle("/report.pdf",
		http.StripPrefix("/", http.FileServer(http.Dir("output"))))

	return http.ListenAndServe(":8080", nil)
}
