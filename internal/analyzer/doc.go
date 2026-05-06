// Package analyzer implements static analysis passes for cron-lint.
//
// # Overview
//
// The analyzer package operates on a set of [Job] values, each of which
// pairs a human-readable name with a fully parsed [parser.Schedule].
//
// # Loading jobs
//
// [LoadJobs] reads a simple text format where every non-blank,
// non-comment line contains a job name followed by the five standard
// cron fields:
//
//	backup    0  2  *  *  *
//	report   30  6  *  *  1-5
//
// # Overlap detection
//
// [DetectOverlaps] compares every pair of jobs and reports an
// [OverlapWarning] whenever all five cron fields share at least one
// common value — meaning both jobs would fire at the same instant.
//
// Typical usage:
//
//	jobs, err := analyzer.LoadJobs(f)
//	if err != nil { ... }
//	for _, w := range analyzer.DetectOverlaps(jobs) {
//		fmt.Println(w)
//	}
package analyzer
