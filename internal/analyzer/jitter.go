package analyzer

import (
	"fmt"
	"strings"
)

// JitterWarning describes a job whose schedule would benefit from jitter
// (a randomised offset) to avoid thundering-herd problems.
type JitterWarning struct {
	Job     Job
	Reason  string
	Suggest string
}

// CheckJitter inspects every job and warns when multiple jobs share an
// identical minute-of-hour value, which causes them to all fire at the same
// wall-clock instant and can overload downstream systems.
//
// It also warns about any single job that fires on the exact hour boundary
// (:00) more than once per day, since that is the most common thundering-herd
// pattern in practice.
func CheckJitter(jobs []Job) []JitterWarning {
	var warnings []JitterWarning

	// Group jobs by their sorted minute set.
	type entry struct {
		minutes []int
		jobs    []Job
	}
	groups := map[string]*entry{}

	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}
		key := minuteKey(j.Schedule.Minutes)
		if _, ok := groups[key]; !ok {
			groups[key] = &entry{minutes: j.Schedule.Minutes}
		}
		groups[key].jobs = append(groups[key].jobs, j)
	}

	for _, g := range groups {
		if len(g.jobs) < 2 {
			continue
		}
		names := jobNames(g.jobs)
		for _, j := range g.jobs {
			warnings = append(warnings, JitterWarning{
				Job:    j,
				Reason: fmt.Sprintf("shares fire-time with %s; consider adding jitter", strings.Join(names, ", ")),
				Suggest: suggestJitterMinute(j.Schedule.Minutes),
			})
		}
	}

	return warnings
}

// minuteKey returns a stable string key for a sorted slice of minute values.
func minuteKey(mins []int) string {
	parts := make([]string, len(mins))
	for i, m := range mins {
		parts[i] = fmt.Sprintf("%d", m)
	}
	return strings.Join(parts, ",")
}

// jobNames returns the Name field of each job, falling back to the raw
// expression when the name is empty.
func jobNames(jobs []Job) []string {
	out := make([]string, 0, len(jobs))
	for _, j := range jobs {
		if j.Name != "" {
			out = append(out, j.Name)
		} else {
			out = append(out, j.Raw)
		}
	}
	return out
}

// suggestJitterMinute proposes an alternative minute that is offset by a small
// prime number from the first minute in the set.
func suggestJitterMinute(mins []int) string {
	if len(mins) == 0 {
		return ""
	}
	offset := (mins[0] + 7) % 60
	return fmt.Sprintf("consider shifting the minute field by a few minutes, e.g. %d", offset)
}
