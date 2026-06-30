package ai

import (
	"fmt"
	"strings"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

type Analysis struct {
	ExecutiveSummary string
	Risk             string
	RiskScore        float64
	Recommendations  []string
	TopRisks         []string
}

func Analyze(scan *models.Scan) Analysis {

	analysis := Analysis{}

	alive := 0
	openPorts := 0
	secrets := 0
	findings := 0

	for _, host := range scan.Hosts {

		if host.Alive {
			alive++
		}

		openPorts += len(host.Ports)
		secrets += len(host.Secrets)
		findings += len(host.Findings)
	}

	analysis.ExecutiveSummary = fmt.Sprintf(
		"RedRecon scanned %d hosts. %d hosts are alive with %d open ports detected.",
		len(scan.Hosts),
		alive,
		openPorts,
	)

	score := 0

	if alive > 0 {
		score += 2
	}

	if openPorts > 10 {
		score += 2
	}

	if secrets > 0 {
		score += 3
	}

	if findings > 0 {
		score += 3
	}

	switch {
	case score >= 8:
		analysis.Risk = "Critical"
	case score >= 5:
		analysis.Risk = "High"
	case score >= 3:
		analysis.Risk = "Medium"
	default:
		analysis.Risk = "Low"
	}
	analysis.RiskScore = float64(score)
	if alive > 0 {
		analysis.Recommendations = append(
			analysis.Recommendations,
			"Review all live assets and remove unused services.",
		)
	}

	if openPorts > 10 {
		analysis.Recommendations = append(
			analysis.Recommendations,
			"Review exposed ports and close unnecessary services.",
		)
	}

	if secrets > 0 {
		analysis.Recommendations = append(
			analysis.Recommendations,
			"Rotate exposed credentials immediately.",
		)
	}

	if findings > 0 {
		analysis.Recommendations = append(
			analysis.Recommendations,
			"Investigate all discovered findings.",
		)
	}

	if len(analysis.Recommendations) == 0 {
		analysis.Recommendations = append(
			analysis.Recommendations,
			"No critical issues detected.",
		)
	}

	analysis.ExecutiveSummary = strings.TrimSpace(analysis.ExecutiveSummary)

	return analysis
}
