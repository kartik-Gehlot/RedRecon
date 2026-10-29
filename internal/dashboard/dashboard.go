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
		if decodeErr := json.NewDecoder(file).Decode(&scan); decodeErr != nil {
			file.Close()
			return decodeErr
		}
		file.Close()
	} else if !os.IsNotExist(err) {
		return err
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	mux.Handle("/scan.json", http.StripPrefix(
		"/", http.FileServer(http.Dir("output")),
	))
	mux.Handle("/report.html", http.StripPrefix(
		"/", http.FileServer(http.Dir("output")),
	))
	mux.Handle("/report.pdf", http.StripPrefix(
		"/", http.FileServer(http.Dir("output")),
	))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Dashboard listening on port " + port)
	return http.ListenAndServe(":"+port, mux)
}
