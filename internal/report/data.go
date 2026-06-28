package report

import "github.com/kartik-Gehlot/RedRecon/internal/models"

type ReportData struct {
	Target       string
	Version      string
	Status       string
	RiskScore    float64
	HostCount    int
	AliveCount   int
	SecretCount  int
	FindingCount int

	Hosts []models.Host
}

func Build(scan *models.Scan) ReportData {

	data := ReportData{
		Target:    scan.Target,
		Version:   scan.Version,
		Status:    scan.Status,
		RiskScore: scan.RiskScore,
		HostCount: len(scan.Hosts),
		Hosts:     scan.Hosts,
	}

	for _, host := range scan.Hosts {

		if host.Alive {
			data.AliveCount++
		}

		data.SecretCount += len(host.Secrets)
		data.FindingCount += len(host.Findings)
	}

	return data
}
