package analyzer

import (
	"fmt"
	"sort"

	"github.com/your-org/cron-lint/internal/parser"
)

// SkewWarning describes a clock-skew risk where multiple jobs share the same
// hour:minute pattern across different weekdays or months, causing uneven load
// distribution when the schedule drifts across DST boundaries.
type SkewWarning struct {
	Jobs    []string
	Minute  int
	Hour    int
	Pattern string
	Message string
}

// CheckSkew detects jobs that fire at the same clock time but on different
// subsets of days, which can cause skewed load after daylight-saving transitions.
func CheckSkew(jobs []Job) []SkewWarning {
	type key struct{ hour, minute int }
	groups := make(map[key][]string)

	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}
		for _, h := range j.Schedule.Hours {
			for _, m := range j.Schedule.Minutes {
				k := key{h, m}
				groups[k] = append(groups[k], j.Name)
			}
		}
	}

	var warnings []SkewWarning
	for k, names := range groups {
		if len(names) < 2 {
			continue
		}
		sort.Strings(names)
		warnings = append(warnings, SkewWarning{
			Jobs:    names,
			Minute:  k.minute,
			Hour:    k.hour,
			Pattern: fmt.Sprintf("%02d:%02d", k.hour, k.minute),
			Message: fmt.Sprintf(
				"%d jobs share clock time %02d:%02d — consider spreading minutes to reduce skew",
				len(names), k.hour, k.minute,
			),
		})
	}

	sort.Slice(warnings, func(i, j int) bool {
		if warnings[i].Hour != warnings[j].Hour {
			return warnings[i].Hour < warnings[j].Hour
		}
		return warnings[i].Minute < warnings[j].Minute
	})
	return warnings
}

// FormatSkewWarnings returns human-readable lines for the given warnings.
func FormatSkewWarnings(ws []SkewWarning) []string {
	lines := make([]string, 0, len(ws))
	for _, w := range ws {
		lines = append(lines, fmt.Sprintf("[SKEW] %s: %v", w.Pattern, w.Message))
	}
	return lines
}

// ensure parser import is used transitively via Job.Schedule.
var _ = parser.Schedule{}
