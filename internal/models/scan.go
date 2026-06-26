package models

type Scan struct {
	Target      string
	Status      string
	Version     string
	StartTime   string
	EndTime     string
	RiskScore   float64
	Subdomains  []string
}