package dashboard

import "github.com/kartik-Gehlot/RedRecon/internal/models"

type DashboardData struct {
	Target      string
	Version     string
	Status      string
	RiskScore   float64

	HostCount   int
	AliveCount  int
	DeadCount   int
	PortCount   int
	SecretCount int

	Hosts []models.Host
}

func Build(scan *models.Scan) DashboardData {

	data := DashboardData{
		Target:    scan.Target,
		Version:   scan.Version,
		Status:    scan.Status,
		RiskScore: scan.RiskScore,
		Hosts:     scan.Hosts,
	}

	for _, host := range scan.Hosts {

		data.HostCount++

		if host.Alive {
			data.AliveCount++
		} else {
			data.DeadCount++
		}

		data.PortCount += len(host.Ports)
		data.SecretCount += len(host.Secrets)
	}

	return data
}