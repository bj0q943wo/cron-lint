// Package analyzer provides static analysis passes for cron job schedules.
//
// # Passes
//
// DetectDuplicates identifies jobs that share an identical normalised schedule
// expression, which is often a copy-paste mistake.
//
// DetectOverlaps finds pairs of jobs whose expanded minute sets intersect,
// meaning they would fire at the same wall-clock minute.
//
// ValidateJobs checks each schedule for semantic problems such as unreachable
// day-of-month values (e.g. the 31st in February).
//
// SuggestFixes inspects schedules for common anti-patterns — such as running
// every minute or clustering many jobs at midnight — and proposes less
// aggressive alternatives.
//
// # Usage
//
//	jobs, err := analyzer.LoadJobs(r)
//	if err != nil { ... }
//
//	duplicates  := analyzer.DetectDuplicates(jobs)
//	overlaps    := analyzer.DetectOverlaps(jobs)
//	warnings    := analyzer.ValidateJobs(jobs)
//	suggestions := analyzer.SuggestFixes(jobs)
package analyzer
