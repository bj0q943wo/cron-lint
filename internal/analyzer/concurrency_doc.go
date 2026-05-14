// Package analyzer provides static-analysis passes for cron schedules.
//
// # Concurrency Check
//
// CheckConcurrency identifies groups of jobs that are scheduled to fire at
// the exact same minute within the week. Concurrent job starts can cause:
//
//   - Resource contention (database connections, file locks, CPU)
//   - Race conditions when jobs share mutable state
//   - Thundering-herd effects on downstream services
//
// The check operates on the full week grid (7 × 24 × 60 = 10 080 minute
// slots). Each job's expanded schedule is mapped onto that grid; any slot
// occupied by two or more jobs yields a ConcurrencyWarning.
//
// Remediation: stagger job start times by at least one minute, or introduce
// a small random jitter at the application layer.
package analyzer
