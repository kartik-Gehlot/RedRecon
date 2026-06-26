package models

type Host struct {
	Hostname     string
	IP           string
	Alive        bool
	Ports        []int
	Technologies []string
}
