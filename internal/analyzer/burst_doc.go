// Package analyzer provides static-analysis passes for cron schedules.
//
// # Burst Detection
//
// CheckBursts scans all jobs and identifies time windows in which an unusually
// large number of jobs are scheduled to run simultaneously (or near-simultaneously).
//
// A "burst" is defined as threshold-or-more jobs firing within a rolling window
// of windowMinutes minutes starting at the same minute-of-hour.
//
// Example usage:
//
//	warnings := analyzer.CheckBursts(jobs, 5, 3)
//	for _, w := range warnings {
//	    fmt.Println(w)
//	}
//
// Overlapping windows are collapsed so that each minute is reported at most once.
package analyzer
