// Package analyzer provides static analysis passes for cron job schedules.
//
// # Timezone Analysis
//
// The timezone sub-feature detects jobs whose scheduled hours fall inside
// Daylight Saving Time (DST) transition windows for a given IANA location:
//
//   - Spring-forward transitions skip one hour entirely, so a job scheduled
//     for that hour will never run on transition day.
//
//   - Fall-back transitions repeat one hour, so a job scheduled for that hour
//     may run twice on transition day, potentially causing duplicate side-effects.
//
// Usage:
//
//	warnings := analyzer.CheckTimezones(jobs, loc)
//	for _, w := range warnings {
//		fmt.Printf("[timezone] %s: %s\n", w.Job.Name, w.Message)
//	}
//
// Pass nil as the location to default to UTC (no DST, no warnings).
package analyzer
