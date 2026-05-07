// Package analyzer provides static analysis passes for cron schedules.
//
// # Frequency Analysis
//
// AnalyzeFrequency inspects each job's parsed schedule and computes:
//
//   - RunsPerHour  – distinct minute slots within one hour
//   - RunsPerDay   – RunsPerHour × distinct hour slots in one day
//   - Category     – "high" (≥30/h), "medium" (≥6/h), or "low" (<6/h)
//
// These metrics help operators identify unexpectedly noisy schedules before
// they reach production.  The reporter package surfaces high-frequency jobs
// as warnings alongside overlap and duplicate diagnostics.
package analyzer
