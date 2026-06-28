package models

import "time"

type Scan struct {
	Target         string
	Status         string
	Version        string
	StartTime      time.Time
	EndTime        time.Time
	RiskScore      float64
	GithubFindings []GitHubFinding
	Hosts          []Host
}
