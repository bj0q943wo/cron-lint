package analyzer

import (
	"fmt"
	"strings"

	"github.com/cron-lint/internal/parser"
)

// Job represents a named cron job with its parsed schedule.
type Job struct {
	Name       string
	Expression string
	Schedule   *parser.Schedule
	Location   interface{} // *time.Location, kept as interface to avoid import cycle
}

// DependencyWarning describes a scheduling dependency concern between two jobs.
type DependencyWarning struct {
	JobA    string
	JobB    string
	Kind    string
	Message string
}

// CheckDependencies inspects a list of jobs for scheduling patterns that suggest
// implicit ordering dependencies: job B always fires within the same minute as
// job A, which may cause race conditions if B depends on A's output.
func CheckDependencies(jobs []Job) []DependencyWarning {
	var warnings []DependencyWarning

	for i := 0; i < len(jobs); i++ {
		for j := i + 1; j < len(jobs); j++ {
			a, b := jobs[i], jobs[j]
			if a.Schedule == nil || b.Schedule == nil {
				continue
			}

			if schedulesOverlap(a.Schedule, b.Schedule) {
				warnings = append(warnings, DependencyWarning{
					JobA:    a.Name,
					JobB:    b.Name,
					Kind:    "concurrent",
					Message: fmt.Sprintf("jobs %q and %q fire at the same time; if one depends on the other, a race condition may occur", a.Name, b.Name),
				})
				continue
			}

			if w, ok := checkSuccessorPattern(a, b); ok {
				warnings = append(warnings, w)
			}
		}
	}

	return warnings
}

// checkSuccessorPattern detects when job B fires exactly one minute after job A
// on every shared trigger, hinting at an implicit pipeline.
func checkSuccessorPattern(a, b Job) (DependencyWarning, bool) {
	shiftedMins := shiftMinutes(a.Schedule.Minutes, 1)
	overlap := intersect(shiftedMins, b.Schedule.Minutes)
	if len(overlap) == 0 {
		return DependencyWarning{}, false
	}

	// Hours, DOM, Month, DOW must all share at least one common value.
	if len(intersect(a.Schedule.Hours, b.Schedule.Hours)) == 0 {
		return DependencyWarning{}, false
	}

	return DependencyWarning{
		JobA:    a.Name,
		JobB:    b.Name,
		Kind:    "successor",
		Message: fmt.Sprintf("job %q fires 1 minute before %q on shared hours [%s]; consider an explicit orchestration dependency", a.Name, b.Name, formatInts(intersect(a.Schedule.Hours, b.Schedule.Hours))),
	}, true
}

func shiftMinutes(mins []int, by int) []int {
	out := make([]int, len(mins))
	for i, m := range mins {
		out[i] = (m + by) % 60
	}
	return out
}

func formatInts(vals []int) string {
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, ",")
}
