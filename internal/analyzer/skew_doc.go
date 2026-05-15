// Package analyzer — skew.go
//
// CheckSkew identifies jobs that are scheduled at the same wall-clock time
// (hour + minute) but on different subsets of days or months.  When a host
// crosses a DST boundary the effective UTC offset changes, causing all of
// those jobs to fire simultaneously instead of being spread across the day.
//
// Example problematic pattern:
//
//	"0 9 * * 1"  – every Monday at 09:00
//	"0 9 * * 3"  – every Wednesday at 09:00
//
// Both jobs share the 09:00 slot.  After a DST spring-forward they will
// both execute at the same moment, potentially overloading shared resources.
//
// Recommendation: offset each job by at least one minute, e.g. 09:00, 09:01.
package analyzer
