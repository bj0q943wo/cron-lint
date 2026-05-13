package analyzer

import (
	"fmt"
	"sort"

	"github.com/your-org/cron-lint/internal/parser"
)

// DriftWarning describes a job whose schedule drifts significantly relative
// to a reference job — e.g. two jobs intended to run together but offset by
// an unexpected number of minutes.
type DriftWarning struct {
	JobA      string
	JobB      string
	OffsetMin int
	Message   string
}

// CheckDrift detects pairs of jobs that share the same hour/day pattern but
// whose minute sets are offset by more than maxOffsetMin minutes apart on
// every shared firing opportunity. A typical use-case: two jobs that should
// run at the same time but have accidentally diverged.
func CheckDrift(jobs []Job, maxOffsetMin int) []DriftWarning {
	var warnings []DriftWarning

	for i := 0; i < len(jobs); i++ {
		for j := i + 1; j < len(jobs); j++ {
			a, b := jobs[i], jobs[j]
			if a.Schedule == nil || b.Schedule == nil {
				continue
			}
			if !sameNonMinutePattern(a.Schedule, b.Schedule) {
				continue
			}
			offset := minSetOffset(a.Schedule.Minutes, b.Schedule.Minutes)
			if offset > maxOffsetMin {
				warnings = append(warnings, DriftWarning{
					JobA:      a.Name,
					JobB:      b.Name,
					OffsetMin: offset,
					Message: fmt.Sprintf(
						"jobs %q and %q share the same hour/day pattern but minute sets differ by %d min (threshold: %d)",
						a.Name, b.Name, offset, maxOffsetMin,
					),
				})
			}
		}
	}
	return warnings
}

// sameNonMinutePattern returns true when two schedules have identical
// hours, days-of-month, months, and days-of-week fields.
func sameNonMinutePattern(a, b *parser.Schedule) bool {
	return intSliceEqual(a.Hours, b.Hours) &&
		intSliceEqual(a.DaysOfMonth, b.DaysOfMonth) &&
		intSliceEqual(a.Months, b.Months) &&
		intSliceEqual(a.DaysOfWeek, b.DaysOfWeek)
}

// minSetOffset returns the minimum absolute difference between any element
// of set a and any element of set b.
func minSetOffset(a, b []int) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	sortedA := sortedCopy(a)
	sortedB := sortedCopy(b)
	min := abs(sortedA[0] - sortedB[0])
	for _, va := range sortedA {
		for _, vb := range sortedB {
			if d := abs(va - vb); d < min {
				min = d
			}
		}
	}
	return min
}

func sortedCopy(s []int) []int {
	c := make([]int, len(s))
	copy(c, s)
	sort.Ints(c)
	return c
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	sa, sb := sortedCopy(a), sortedCopy(b)
	for i := range sa {
		if sa[i] != sb[i] {
			return false
		}
	}
	return true
}
