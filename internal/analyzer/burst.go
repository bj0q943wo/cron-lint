package analyzer

import (
	"fmt"
	"strings"
)

// BurstWarning describes a group of jobs that fire within a short time window,
// potentially causing a resource burst.
type BurstWarning struct {
	// WindowMinutes is the rolling window size used for detection.
	WindowMinutes int
	// Jobs contains the names/expressions of jobs in the burst group.
	Jobs []JobEntry
	// Minute is the minute-of-hour where the burst begins.
	Minute int
	// Count is how many jobs fire within the window starting at Minute.
	Count int
}

func (w BurstWarning) String() string {
	names := make([]string, len(w.Jobs))
	for i, j := range w.Jobs {
		names[i] = j.Name
	}
	return fmt.Sprintf(
		"burst of %d jobs starting at minute %02d (within %d-minute window): %s",
		w.Count, w.Minute, w.WindowMinutes, strings.Join(names, ", "),
	)
}

// CheckBursts detects groups of jobs that all fire within a rolling window of
// windowMinutes minutes in the same hour, which may indicate a resource spike.
// Only hours where at least threshold jobs overlap are reported.
func CheckBursts(jobs []JobEntry, windowMinutes, threshold int) []BurstWarning {
	if windowMinutes <= 0 || threshold <= 0 {
		return nil
	}

	var warnings []BurstWarning
	reported := make(map[int]bool)

	for startMinute := 0; startMinute < 60; startMinute++ {
		if reported[startMinute] {
			continue
		}
		var burst []JobEntry
		for _, job := range jobs {
			if job.Schedule == nil {
				continue
			}
			if firesInWindow(job, startMinute, windowMinutes) {
				burst = append(burst, job)
			}
		}
		if len(burst) >= threshold {
			warnings = append(warnings, BurstWarning{
				WindowMinutes: windowMinutes,
				Jobs:          burst,
				Minute:        startMinute,
				Count:         len(burst),
			})
			for m := startMinute; m < startMinute+windowMinutes && m < 60; m++ {
				reported[m] = true
			}
		}
	}
	return warnings
}

// firesInWindow returns true if the job fires at least once during
// [startMinute, startMinute+windowMinutes) in any hour.
func firesInWindow(job JobEntry, startMinute, windowMinutes int) bool {
	for m := startMinute; m < startMinute+windowMinutes && m < 60; m++ {
		for _, jm := range job.Schedule.Minutes {
			if jm == m {
				return true
			}
		}
	}
	return false
}
