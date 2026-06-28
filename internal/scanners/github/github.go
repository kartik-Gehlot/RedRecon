package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/kartik-Gehlot/RedRecon/internal/models"
)

func Run(scan *models.Scan) error {

	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		fmt.Println("[WARN] GITHUB_TOKEN not configured. Skipping GitHub Search.")
		return nil
	}

	query := url.QueryEscape(scan.Target)

	apiURL := "https://api.github.com/search/code?q=" + query

	req, err := http.NewRequest(
		"GET",
		apiURL,
		nil,
	)

	if err != nil {
		return err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api returned status %d", resp.StatusCode)
	}
	var result SearchResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}
	for _, item := range result.Items {

		scan.GithubFindings = append(
			scan.GithubFindings,
			models.GitHubFinding{
				Repository: item.Repository.FullName,
				File:       item.Path,
				URL:        item.HTMLURL,
			},
		)
	}

	return nil
}
