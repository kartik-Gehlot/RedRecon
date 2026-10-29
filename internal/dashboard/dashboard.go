package dashboard

import (
	"crypto/subtle"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/kartik-Gehlot/RedRecon/internal/ai"
	"github.com/kartik-Gehlot/RedRecon/internal/logger"
	"github.com/kartik-Gehlot/RedRecon/internal/models"
	"github.com/kartik-Gehlot/RedRecon/internal/orchestrator"
)

type DashboardPage struct {
	DashboardData
	AI ai.Analysis
}

var scanState = struct {
	sync.Mutex
	Running bool   `json:"running"`
	Target  string `json:"target"`
	Status  string `json:"status"`
	Error   string `json:"error"`
}{
	Status: "Idle",
}

func tokenMatches(given, expected string) bool {
	return given != "" &&
		subtle.ConstantTimeCompare([]byte(given), []byte(expected)) == 1
}

func Start() error {
	if err := os.MkdirAll("output", 0755); err != nil {
		return err
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		var scan models.Scan

		file, err := os.Open("output/scan.json")
		if err == nil {
			decodeErr := json.NewDecoder(file).Decode(&scan)
			file.Close()

			if decodeErr != nil {
				http.Error(w, "Unable to read scan results", http.StatusInternalServerError)
				return
			}
		} else if !os.IsNotExist(err) {
			http.Error(w, "Unable to open scan results", http.StatusInternalServerError)
			return
		}

		analysis := ai.Analyze(&scan)
		scan.RiskScore = analysis.RiskScore

		page := DashboardPage{
			DashboardData: Build(&scan),
			AI:            analysis,
		}

		if err := tmpl.Execute(w, page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/api/scan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		expectedToken := os.Getenv("REDRECON_SCAN_TOKEN")
		if expectedToken == "" ||
			!tokenMatches(r.Header.Get("X-Scan-Token"), expectedToken) {
			http.Error(w, "Unauthorized or scanning not configured", http.StatusUnauthorized)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		defer r.Body.Close()

		var request struct {
			Target string `json:"target"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		target := strings.TrimSpace(request.Target)
		if target == "" || strings.ContainsAny(target, "/:@ \t\r\n") {
			http.Error(w, "Provide a valid domain name", http.StatusBadRequest)
			return
		}

		scanState.Lock()
		if scanState.Running {
			scanState.Unlock()
			http.Error(w, "A scan is already running", http.StatusConflict)
			return
		}

		scanState.Running = true
		scanState.Target = target
		scanState.Status = "Running"
		scanState.Error = ""
		scanState.Unlock()

		go func() {
			err := orchestrator.Start(target)

			scanState.Lock()
			defer scanState.Unlock()

			scanState.Running = false
			if err != nil {
				scanState.Status = "Failed"
				scanState.Error = err.Error()
			} else {
				scanState.Status = "Completed"
			}
		}()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "Running",
			"target": target,
		})
	})

	mux.HandleFunc("/api/scan/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		expectedToken := os.Getenv("REDRECON_SCAN_TOKEN")
		if expectedToken == "" ||
			!tokenMatches(r.Header.Get("X-Scan-Token"), expectedToken) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		scanState.Lock()
		defer scanState.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(scanState)
	})

	for _, path := range []string{"/scan.json", "/report.html", "/report.pdf"} {
		mux.Handle(path, http.StripPrefix("/", http.FileServer(http.Dir("output"))))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Dashboard listening on port " + port)
	return http.ListenAndServe(":"+port, mux)
}
