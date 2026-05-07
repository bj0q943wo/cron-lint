package analyzer

import (
	"fmt"

	"github.com/user/cron-lint/internal/parser"
)

// FrequencyReport summarises how often a job fires.
type FrequencyReport struct {
	Job            Job
	RunsPerHour    int
	RunsPerDay     int
	Category       string // "high", "medium", "low"
}

// AnalyzeFrequency computes execution frequency for each job and
// categorises it so downstream reporters can surface noisy schedules.
func AnalyzeFrequency(jobs []Job) []FrequencyReport {
	reports := make([]FrequencyReport, 0, len(jobs))
	for _, j := range jobs {
		rph := runsPerHour(j.Schedule)
		rpd := rph * hoursPerDay(j.Schedule)
		reports = append(reports, FrequencyReport{
			Job:         j,
			RunsPerHour: rph,
			RunsPerDay:  rpd,
			Category:    category(rph),
		})
	}
	return reports
}

// runsPerHour returns the number of distinct minutes within a single hour
// that the schedule fires on (ignoring day/month/weekday constraints).
func runsPerHour(s parser.Schedule) int {
	return len(s.Minutes)
}

// hoursPerDay returns the number of distinct hours in a day the schedule fires.
func hoursPerDay(s parser.Schedule) int {
	return len(s.Hours)
}

func category(rph int) string {
	switch {
	case rph >= 30:
		return "high"
	case rph >= 6:
		return "medium"
	default:
		return "low"
	}
}

// FormatFrequency returns a human-readable frequency string.
func FormatFrequency(r FrequencyReport) string {
	return fmt.Sprintf("%s: ~%d run(s)/hour, ~%d run(s)/day [%s]",
		r.Job.Name, r.RunsPerHour, r.RunsPerDay, r.Category)
}
