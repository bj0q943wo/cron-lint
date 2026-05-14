package analyzer

import "fmt"

// ConcurrencyWarning describes a set of jobs that may run concurrently
// within the same minute slot, potentially competing for shared resources.
type ConcurrencyWarning struct {
	// Minute is the absolute minute-of-week (0..10079) where the collision occurs.
	Minute int
	// Jobs holds the names of the jobs that all fire at that minute.
	Jobs []string
	// Suggestion is a human-readable remediation hint.
	Suggestion string
}

// CheckConcurrency detects groups of jobs that fire in the same minute slot
// across the full week (0..10079 minutes). Any slot occupied by two or more
// jobs produces a ConcurrencyWarning.
//
// Jobs with a nil Schedule are silently skipped.
func CheckConcurrency(jobs []Job) []ConcurrencyWarning {
	// minute-of-week -> list of job names
	slots := make(map[int][]string)

	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}
		seen := make(map[int]bool)
		for _, dow := range j.Schedule.DayOfWeek {
			for _, h := range j.Schedule.Hours {
				for _, m := range j.Schedule.Minutes {
					slot := dow*24*60 + h*60 + m
					if !seen[slot] {
						slots[slot] = append(slots[slot], j.Name)
						seen[slot] = true
					}
				}
			}
		}
	}

	var warnings []ConcurrencyWarning
	for minute, names := range slots {
		if len(names) < 2 {
			continue
		}
		warnings = append(warnings, ConcurrencyWarning{
			Minute:     minute,
			Jobs:       names,
			Suggestion: fmt.Sprintf("Stagger start times for: %v", names),
		})
	}
	return warnings
}
