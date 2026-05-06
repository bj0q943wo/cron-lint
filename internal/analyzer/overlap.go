// Package analyzer provides static analysis for cron expressions,
// including overlap detection between scheduled jobs.
package analyzer

import (
	"fmt"

	"github.com/cron-lint/internal/parser"
)

// Job represents a named cron job with its parsed schedule.
type Job struct {
	Name       string
	Expression string
	Schedule   *parser.Schedule
}

// OverlapWarning describes two jobs whose schedules intersect.
type OverlapWarning struct {
	JobA    string
	JobB    string
	Minutes []int
	Hours   []int
}

func (w OverlapWarning) String() string {
	return fmt.Sprintf("overlap detected between %q and %q", w.JobA, w.JobB)
}

// DetectOverlaps returns all pairs of jobs that share at least one
// (minute, hour, dom, month, dow) intersection.
func DetectOverlaps(jobs []Job) []OverlapWarning {
	var warnings []OverlapWarning
	for i := 0; i < len(jobs); i++ {
		for j := i + 1; j < len(jobs); j++ {
			a, b := jobs[i], jobs[j]
			if schedulesOverlap(a.Schedule, b.Schedule) {
				warnings = append(warnings, OverlapWarning{
					JobA:    a.Name,
					JobB:    b.Name,
					Minutes: intersect(a.Schedule.Minutes, b.Schedule.Minutes),
					Hours:   intersect(a.Schedule.Hours, b.Schedule.Hours),
				})
			}
		}
	}
	return warnings
}

// schedulesOverlap returns true when every field has at least one common value.
func schedulesOverlap(a, b *parser.Schedule) bool {
	return len(intersect(a.Minutes, b.Minutes)) > 0 &&
		len(intersect(a.Hours, b.Hours)) > 0 &&
		len(intersect(a.DaysOfMonth, b.DaysOfMonth)) > 0 &&
		len(intersect(a.Months, b.Months)) > 0 &&
		len(intersect(a.DaysOfWeek, b.DaysOfWeek)) > 0
}

// intersect returns the sorted values present in both slices.
func intersect(a, b []int) []int {
	set := make(map[int]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}
	var result []int
	for _, v := range a {
		if _, ok := set[v]; ok {
			result = append(result, v)
		}
	}
	return result
}
