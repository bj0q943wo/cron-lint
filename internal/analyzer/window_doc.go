// Package analyzer — window analysis
//
// window.go provides CheckWindows, which identifies cron jobs whose entire
// schedule falls within a user-defined sensitive time window such as a nightly
// maintenance period or a peak-traffic window.
//
// Usage:
//
//	windows := []analyzer.TimeWindow{
//		{Name: "maintenance", Start: 2 * 60, End: 4 * 60},
//	}
//	warnings := analyzer.CheckWindows(jobs, windows)
//	fmt.Println(analyzer.FormatWindowWarnings(warnings))
//
// A warning is emitted only when *every* (hour, minute) pair produced by the
// schedule lies within the window — jobs that sometimes fire outside the window
// are not flagged.
package analyzer
