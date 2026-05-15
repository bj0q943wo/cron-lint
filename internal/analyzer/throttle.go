package analyzer

import (
	"fmt"

	"github.com/your-org/cron-lint/internal/parser"
)

// ThrottleOptions configures the throttle checker.
type ThrottleOptions struct {
	// MaxFiringsPer5Min is the maximum allowed job firings within any 5-minute window.
	MaxFiringsPer5Min int
}

// DefaultThrottleOptions returns sensible defaults.
var DefaultThrottleOptions = ThrottleOptions{
	MaxFiringsPer5Min: 10,
}

// ThrottleWarning is emitted when a group of jobs fires too frequently in a
// short window, risking resource exhaustion or rate-limit violations.
type ThrottleWarning struct {
	Window    string   // human-readable window label, e.g. "00:00-00:05"
	Jobs      []string // job names that fire in the window
	Firings   int      // total firings across all jobs
	Threshold int      // configured threshold
}

func (w ThrottleWarning) String() string {
	return fmt.Sprintf(
		"window %s: %d firings across %d jobs exceeds threshold of %d",
		w.Window, w.Firings, len(w.Jobs), w.Threshold,
	)
}

// CheckThrottle detects 5-minute windows where the combined firing count of
// all provided jobs exceeds opts.MaxFiringsPer5Min.
func CheckThrottle(jobs []Job, opts ThrottleOptions) []ThrottleWarning {
	type slot struct {
		names    map[string]struct{}
		firings  int
	}

	// slots keyed by (hour*12 + 5-min-block)
	slots := make(map[int]*slot)

	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}
		for _, h := range j.Schedule.Hours {
			for _, m := range j.Schedule.Minutes {
				key := h*12 + m/5
				if _, ok := slots[key]; !ok {
					slots[key] = &slot{names: make(map[string]struct{})}
				}
				slots[key].names[j.Name] = struct{}{}
				slots[key].firings++
			}
		}
	}

	var warnings []ThrottleWarning
	for key, s := range slots {
		if s.firings <= opts.MaxFiringsPer5Min {
			continue
		}
		h := key / 12
		block := (key % 12) * 5
		win := fmt.Sprintf("%02d:%02d-%02d:%02d", h, block, h, block+4)
		names := make([]string, 0, len(s.names))
		for n := range s.names {
			names = append(names, n)
		}
		warnings = append(warnings, ThrottleWarning{
			Window:    win,
			Jobs:      names,
			Firings:   s.firings,
			Threshold: opts.MaxFiringsPer5Min,
		})
	}
	return warnings
}

// firesInSlot returns how many times schedule s fires within the given
// 5-minute block (blockStart .. blockStart+4 inclusive) for every hour.
func firesInSlot(s *parser.Schedule, blockStart int) int {
	count := 0
	for _, m := range s.Minutes {
		if m >= blockStart && m <= blockStart+4 {
			count += len(s.Hours)
		}
	}
	return count
}
