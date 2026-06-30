package github

type SearchResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name       string     `json:"name"`
	HTMLURL    string     `json:"html_url"`
	Path       string     `json:"path"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	FullName string `json:"full_name"`
}
