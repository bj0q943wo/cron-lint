package analyzer

import (
	"fmt"
	"time"

	"github.com/user/cron-lint/internal/parser"
)

// StalenessWarning describes a job whose schedule has not fired recently
// or will not fire within a reasonable future window.
type StalenessWarning struct {
	JobName    string
	Expression string
	Message    string
}

// CheckStaleness inspects each job and warns when its schedule would not
// fire within the next staleDays calendar days, suggesting the cron
// expression may be misconfigured or intentionally infrequent.
//
// A window of 366 days is used so that yearly jobs (e.g. "0 9 1 1 *")
// are not flagged, while expressions that can never fire (e.g. Feb 30)
// are caught.
func CheckStaleness(jobs []Job, reference time.Time) []StalenessWarning {
	const lookAheadDays = 366

	var warnings []StalenessWarning
	deadline := reference.Add(time.Duration(lookAheadDays) * 24 * time.Hour)

	for _, job := range jobs {
		if job.Schedule == nil {
			continue
		}
		if !firesBeforeDeadline(job.Schedule, reference, deadline) {
			warnings = append(warnings, StalenessWarning{
				JobName:    job.Name,
				Expression: job.Expression,
				Message: fmt.Sprintf(
					"schedule does not fire within the next %d days; verify the expression is correct",
					lookAheadDays,
				),
			})
		}
	}
	return warnings
}

// firesBeforeDeadline returns true when the schedule has at least one
// minute-level match between reference (inclusive) and deadline (exclusive).
func firesBeforeDeadline(s *parser.Schedule, reference, deadline time.Time) bool {
	// Walk day-by-day to keep the inner loop bounded.
	for d := reference.Truncate(24 * time.Hour); d.Before(deadline); d = d.Add(24 * time.Hour) {
		month := int(d.Month())
		dom := d.Day()
		dow := int(d.Weekday())

		if !contains(s.Months, month) {
			continue
		}
		if !contains(s.DaysOfMonth, dom) {
			continue
		}
		if !contains(s.DaysOfWeek, dow) {
			continue
		}
		// At least one hour+minute combination exists on this day.
		if len(s.Hours) > 0 && len(s.Minutes) > 0 {
			return true
		}
	}
	return false
}

func contains(slice []int, v int) bool {
	for _, x := range slice {
		if x == v {
			return true
		}
	}
	return false
}
