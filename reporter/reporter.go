package reporter

import (
	"context"
	"fmt"

	"github.com/igumus/task-sniffer/config"
)

type ReportFunc func(context.Context, <-chan Task) error

type Task struct {
	Row     int
	Keyword string
	File    string
	Desc    string
}

func DefaultReporter(config config.Config) (ReportFunc, error) {
	return func(ctx context.Context, tasks <-chan Task) error {
		fmt.Printf("[Project Addr]     %s\n", config.Addr())
		fmt.Printf("[Project Name]     %s\n", config.Name())
		fmt.Printf("[Project Reporter] Console\n\n")
		for t := range tasks {
			fmt.Printf(" - [ ] %s: %s\n\t* at %s:%d\n", t.Keyword, t.Desc, t.File, t.Row)
		}
		return nil
	}, nil
}
