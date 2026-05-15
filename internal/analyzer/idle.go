package analyzer

import (
	"fmt"
	"sort"
)

// IdleWarning describes a gap in scheduling where no jobs run for an extended period.
type IdleWarning struct {
	StartHour int
	EndHour   int
	GapHours  int
	Message   string
}

// IdleOptions controls the sensitivity of idle-gap detection.
type IdleOptions struct {
	// MinGapHours is the minimum consecutive hours with no jobs before a warning is raised.
	MinGapHours int
}

// DefaultIdleOptions returns sensible defaults.
func DefaultIdleOptions() IdleOptions {
	return IdleOptions{MinGapHours: 4}
}

// CheckIdle inspects the set of jobs and warns when there are long stretches of
// hours (across a 24-hour day) during which no job is scheduled to run.
func CheckIdle(jobs []Job, opts IdleOptions) []IdleWarning {
	if opts.MinGapHours <= 0 {
		opts.MinGapHours = DefaultIdleOptions().MinGapHours
	}

	active := make(map[int]bool)
	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}
		for _, h := range j.Schedule.Hours {
			active[h] = true
		}
	}

	var warnings []IdleWarning
	gapStart := -1
	gapLen := 0

	check := func(h int) {
		if !active[h%24] {
			if gapStart < 0 {
				gapStart = h % 24
			}
			gapLen++
		} else {
			if gapLen >= opts.MinGapHours {
				end := (gapStart + gapLen) % 24
				warnings = append(warnings, IdleWarning{
					StartHour: gapStart,
					EndHour:   end,
					GapHours:  gapLen,
					Message:   fmt.Sprintf("no jobs scheduled for %d consecutive hours (%02d:00–%02d:00)", gapLen, gapStart, end),
				})
			}
			gapStart = -1
			gapLen = 0
		}
	}

	for h := 0; h < 24; h++ {
		check(h)
	}
	// close any trailing gap
	if gapLen >= opts.MinGapHours {
		end := (gapStart + gapLen) % 24
		warnings = append(warnings, IdleWarning{
			StartHour: gapStart,
			EndHour:   end,
			GapHours:  gapLen,
			Message:   fmt.Sprintf("no jobs scheduled for %d consecutive hours (%02d:00–%02d:00)", gapLen, gapStart, end),
		})
	}

	sort.Slice(warnings, func(i, j int) bool {
		return warnings[i].StartHour < warnings[j].StartHour
	})
	return warnings
}
