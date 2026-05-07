// Package analyzer provides static analysis passes for cron job schedules.
//
// # Dependency Analysis
//
// The dependency sub-feature detects implicit ordering concerns between jobs:
//
//   - Concurrent: two or more jobs are scheduled to fire at the exact same
//     minute. If any of those jobs consumes output produced by another, a race
//     condition exists.
//
//   - Successor: job B is scheduled exactly one minute after job A on the same
//     hour(s). This pattern often indicates an undocumented pipeline where B
//     assumes A has already completed, which is fragile.
//
// Usage:
//
//	warnings := analyzer.CheckDependencies(jobs)
//	for _, w := range warnings {
//		fmt.Println(w.Kind, w.Message)
//	}
package analyzer
