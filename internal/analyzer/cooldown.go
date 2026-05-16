package analyzer

import (
	"fmt"
	"sort"
)

// CooldownWarning is raised when two jobs share overlapping schedules and fire
// too close together, leaving insufficient cooldown time between runs.
type CooldownWarning struct {
	JobA         string
	JobB         string
	MinuteA      int
	MinuteB      int
	GapMinutes   int
	MinRequired  int
}

func (w CooldownWarning) String() string {
	return fmt.Sprintf(
		"jobs %q and %q fire only %d minute(s) apart (minimum required: %d); consider spreading them further",
		w.JobA, w.JobB, w.GapMinutes, w.MinRequired,
	)
}

// CooldownOptions controls the behaviour of CheckCooldown.
type CooldownOptions struct {
	// MinGapMinutes is the minimum number of minutes that must separate any two
	// jobs' fire times within the same hour. Defaults to 5.
	MinGapMinutes int
}

// DefaultCooldownOptions returns sensible defaults.
func DefaultCooldownOptions() CooldownOptions {
	return CooldownOptions{MinGapMinutes: 5}
}

// CheckCooldown inspects all pairs of jobs and warns when two jobs are
// scheduled to fire within fewer than opts.MinGapMinutes of each other within
// any given hour.
func CheckCooldown(jobs []Job, opts CooldownOptions) []CooldownWarning {
	if opts.MinGapMinutes <= 0 {
		opts.MinGapMinutes = DefaultCooldownOptions().MinGapMinutes
	}

	var warnings []CooldownWarning

	for i := 0; i < len(jobs); i++ {
		if jobs[i].Schedule == nil {
			continue
		}
		minutesA := sortedCopy(jobs[i].Schedule.Minutes)

		for j := i + 1; j < len(jobs); j++ {
			if jobs[j].Schedule == nil {
				continue
			}
			minutesB := sortedCopy(jobs[j].Schedule.Minutes)

			if gap, mA, mB, ok := smallestGap(minutesA, minutesB); ok && gap < opts.MinGapMinutes {
				warnings = append(warnings, CooldownWarning{
					JobA:        jobs[i].Name,
					JobB:        jobs[j].Name,
					MinuteA:     mA,
					MinuteB:     mB,
					GapMinutes:  gap,
					MinRequired: opts.MinGapMinutes,
				})
			}
		}
	}
	return warnings
}

// smallestGap returns the smallest absolute minute difference between any
// element of a and any element of b, along with the contributing minutes.
func smallestGap(a, b []int) (gap, mA, mB int, found bool) {
	if len(a) == 0 || len(b) == 0 {
		return 0, 0, 0, false
	}

	sort.Ints(a)
	sort.Ints(b)

	minGap := 61
	bestA, bestB := 0, 0

	ia, ib := 0, 0
	for ia < len(a) && ib < len(b) {
		diff := a[ia] - b[ib]
		if diff < 0 {
			diff = -diff
		}
		if diff < minGap {
			minGap = diff
			bestA, bestB = a[ia], b[ib]
		}
		if a[ia] < b[ib] {
			ia++
		} else {
			ib++
		}
	}
	return minGap, bestA, bestB, true
}
