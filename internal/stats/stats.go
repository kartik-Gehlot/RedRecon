package stats

import (
	"fmt"
	"time"
)

type ModuleStat struct {
	Name     string
	Duration time.Duration
}

var modules []ModuleStat

func Record(name string, start time.Time) {
	modules = append(modules, ModuleStat{
		Name:     name,
		Duration: time.Since(start),
	})
}

func Print() {

	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("        MODULE PERFORMANCE")
	fmt.Println("========================================")

	for _, m := range modules {

		fmt.Printf("%-25s %v\n", m.Name, m.Duration)

	}

	fmt.Println("========================================")
}
