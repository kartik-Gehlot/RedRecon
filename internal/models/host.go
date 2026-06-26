package models

type Host struct {
	Hostname     string
	IP           string
	Alive        bool
	StatusCode   int
	Title        string
	Technologies []string
	Ports        []int
	URLs         []string
	JavaScripts  []string
	Secrets      []string
	Findings     []string
}
