package analyzer

import (
	"fmt"
	"strings"

	"github.com/user/cron-lint/internal/parser"
)

// ValidationWarning represents a non-fatal issue found in a cron expression.
type ValidationWarning struct {
	Job     Job
	Message string
}

// ValidateJobs inspects each job's parsed schedule for suspicious patterns
// and returns a list of warnings. Warnings do not prevent execution but
// indicate schedules that are likely unintentional.
func ValidateJobs(jobs []Job) []ValidationWarning {
	var warnings []ValidationWarning
	for _, job := range jobs {
		warnings = append(warnings, checkSchedule(job)...)
	}
	return warnings
}

func checkSchedule(job Job) []ValidationWarning {
	var ws []ValidationWarning

	// Warn when a job runs every minute (minutes field has all 60 values).
	if len(job.Schedule.Minutes) == 60 {
		ws = append(ws, ValidationWarning{
			Job:     job,
			Message: "schedule fires every minute, which may cause high load",
		})
	}

	// Warn when day-of-week and day-of-month are both fully expanded
	// (all days covered), which is redundant but harmless.
	if len(job.Schedule.DaysOfWeek) == 7 && len(job.Schedule.DaysOfMonth) == 31 {
		ws = append(ws, ValidationWarning{
			Job:     job,
			Message: "both day-of-month and day-of-week are unrestricted; consider constraining one",
		})
	}

	// Warn about jobs that never fire because a specific day-of-month
	// (e.g. 31) combined with months that never have that many days.
	if warns := checkUnreachableDayOfMonth(job); len(warns) > 0 {
		ws = append(ws, warns...)
	}

	return ws
}

// checkUnreachableDayOfMonth warns when day 31 is requested but none of the
// selected months contain 31 days.
func checkUnreachableDayOfMonth(job Job) []ValidationWarning {
	// Months with fewer than 31 days: 2(Feb),4(Apr),6(Jun),9(Sep),11(Nov).
	shortMonths := map[int]bool{2: true, 4: true, 6: true, 9: true, 11: true}

	has31 := false
	for _, d := range job.Schedule.DaysOfMonth {
		if d == 31 {
			has31 = true
			break
		}
	}
	if !has31 {
		return nil
	}

	// Check whether every selected month is a short month.
	allShort := true
	for _, m := range job.Schedule.Months {
		if !shortMonths[m] {
			allShort = false
			break
		}
	}
	if allShort && len(job.Schedule.Months) > 0 {
		return []ValidationWarning{{
			Job:     job,
			Message: fmt.Sprintf("day 31 is specified but month(s) [%s] never have 31 days",
				joinInts(job.Schedule.Months)),
		}}
	}
	return nil
}

func joinInts(ns []int) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = fmt.Sprintf("%d", n)
	}
	return strings.Join(parts, ",")
}

// Ensure parser import is used via the Schedule type.
var _ = parser.Schedule{}
