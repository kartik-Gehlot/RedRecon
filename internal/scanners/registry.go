package scanners

var registry []Scanner

func Register(s Scanner) {
	registry = append(registry, s)
}

func All() []Scanner {
	return registry
}
