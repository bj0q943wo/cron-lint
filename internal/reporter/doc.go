// Package reporter provides formatting and output utilities for cron-lint
// analysis results.
//
// It supports two output formats:
//
//	- text: human-readable, suitable for terminal output
//	- json: machine-readable, suitable for CI pipelines or tooling integration
//
// Typical usage:
//
//	jobs, _ := analyzer.LoadJobs(input)
//	overlaps := analyzer.DetectOverlaps(jobs)
//	report := reporter.Build(jobs, overlaps)
//	reporter.WriteText(os.Stdout, report)
//
// The reporter does not perform any analysis itself; it only formats results
// produced by the analyzer package.
package reporter
