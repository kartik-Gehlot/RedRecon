package dashboard

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/kartik-Gehlot/RedRecon/internal/ai"
	"github.com/kartik-Gehlot/RedRecon/internal/logger"
	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

type DashboardPage struct {
	DashboardData
	AI ai.Analysis
}

//	func Start() error {
//		file, err := os.Open("output/scan.json")
//		if err != nil {
//			return err
//		}
//		defer file.Close()
func Start() error {
	if err := os.MkdirAll("output", 0755); err != nil {
		return err
	}

	var scan models.Scan

file, err := os.Open("output/scan.json")
if err == nil {
    defer file.Close()

    if err := json.NewDecoder(file).Decode(&scan); err != nil {
        return err
    }
} else if !os.IsNotExist(err) {
    return err
}

// Continue loading the template and starting the dashboard.
	defer file.Close()

	// var scan models.Scan

	if err := json.NewDecoder(file).Decode(&scan); err != nil {
		return err
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	http.Handle("/scan.json",
		http.StripPrefix("/", http.FileServer(http.Dir("output"))))

	http.Handle("/report.html",
		http.StripPrefix("/", http.FileServer(http.Dir("output"))))

	http.Handle("/report.pdf",
		http.StripPrefix("/", http.FileServer(http.Dir("output"))))

	// return http.ListenAndServe(":8080", nil)

port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}

logger.Info("Dashboard listening on port " + port)
return http.ListenAndServe(":"+port, nil)

}
