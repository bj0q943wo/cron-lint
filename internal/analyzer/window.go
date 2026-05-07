package analyzer

import (
	"fmt"
	"strings"

	"github.com/your-org/cron-lint/internal/parser"
)

// WindowWarning describes a job whose schedule falls entirely within a narrow
// maintenance or high-load window defined by the caller.
type WindowWarning struct {
	Job     parser.Job
	Window  TimeWindow
	Message string
}

// TimeWindow represents a contiguous range of minutes-of-day [Start, End).
type TimeWindow struct {
	Name  string
	Start int // minutes since midnight, inclusive
	End   int // minutes since midnight, exclusive
}

// CheckWindows reports jobs whose every weekly firing falls inside one of the
// provided sensitive windows (e.g. a nightly maintenance window).
func CheckWindows(jobs []parser.Job, windows []TimeWindow) []WindowWarning {
	var warnings []WindowWarning
	for _, job := range jobs {
		if job.Schedule == nil {
			continue
		}
		for _, w := range windows {
			if allMinutesInWindow(job.Schedule, w) {
				warnings = append(warnings, WindowWarning{
					Job:    job,
					Window: w,
					Message: fmt.Sprintf(
						"job %q always fires inside window %q (%s–%s)",
						job.Name, w.Name,
						minutesToClock(w.Start), minutesToClock(w.End),
					),
				})
				break
			}
		}
	}
	return warnings
}

// allMinutesInWindow returns true when every (hour, minute) combination
// produced by the schedule falls within the window.
func allMinutesInWindow(s *parser.Schedule, w TimeWindow) bool {
	if len(s.Hours) == 0 || len(s.Minutes) == 0 {
		return false
	}
	for _, h := range s.Hours {
		for _, m := range s.Minutes {
			mod := h*60 + m
			if mod < w.Start || mod >= w.End {
				return false
			}
		}
	}
	return true
}

func minutesToClock(min int) string {
	return fmt.Sprintf("%02d:%02d", min/60, min%60)
}

// FormatWindowWarnings returns a human-readable summary of window warnings.
func FormatWindowWarnings(warnings []WindowWarning) string {
	if len(warnings) == 0 {
		return "no window conflicts detected"
	}
	var sb strings.Builder
	for _, w := range warnings {
		sb.WriteString("[WINDOW] ")
		sb.WriteString(w.Message)
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}
