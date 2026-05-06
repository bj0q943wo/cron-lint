// Command cron-lint is a static analyzer for cron expressions.
// It reads a cron-style job file, validates all expressions, and
// reports any scheduling overlaps between jobs.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/cron-lint/internal/analyzer"
	"github.com/example/cron-lint/internal/reporter"
)

func main() {
	var (
		filePath   = flag.String("f", "", "Path to cron job file (required)")
		outputFmt  = flag.String("format", "text", "Output format: text or json")
		exitOnWarn = flag.Bool("strict", false", "Exit with code 1 if any overlaps are found")
	)
	flag.Parse()

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "error: -f <file> is required")
		flag.Usage()
		os.Exit(2)
	}

	f, err := os.Open(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		os.Exit(2)
	}
	defer f.Close()

	jobs, err := analyzer.LoadJobs(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading jobs: %v\n", err)
		os.Exit(2)
	}

	overlaps := analyzer.DetectOverlaps(jobs)
	report := reporter.Build(jobs, overlaps)

	switch *outputFmt {
	case "json":
		if err := reporter.WriteJSON(os.Stdout, report); err != nil {
			fmt.Fprintf(os.Stderr, "error writing JSON: %v\n", err)
			os.Exit(2)
		}
	default:
		reporter.WriteText(os.Stdout, report)
	}

	if *exitOnWarn && len(overlaps) > 0 {
		os.Exit(1)
	}
}
