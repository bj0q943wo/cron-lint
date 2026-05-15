// Package analyzer provides static-analysis passes for cron schedules.
//
// # Idle-gap detection
//
// CheckIdle scans the union of all scheduled jobs across a 24-hour clock and
// reports any consecutive stretch of hours during which no job is configured
// to run.  Long idle windows can indicate misconfigured schedules, forgotten
// maintenance jobs, or monitoring blind-spots.
//
// A warning is emitted when the gap meets or exceeds IdleOptions.MinGapHours
// (default 4).  Each IdleWarning carries the start/end hour and a human-
// readable message suitable for both text and JSON reporters.
package analyzer
