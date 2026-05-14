package analyzer

import (
	"fmt"

	"github.com/your-org/cron-lint/internal/parser"
)

// CardinalityWarning describes a job whose schedule produces an unexpectedly
// large (or small) number of distinct fire-times per week.
type CardinalityWarning struct {
	JobName    string
	Expression string
	FiresPerWeek int
	Message    string
}

// CheckCardinality inspects each job's schedule and warns when the number of
// distinct fire-times per week falls outside the supplied bounds.
//
// lowerBound == 0 disables the low-frequency check.
// upperBound == 0 disables the high-frequency check.
func CheckCardinality(jobs []Job, lowerBound, upperBound int) []CardinalityWarning {
	var warnings []CardinalityWarning

	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}

		count := firesPerWeek(j.Schedule)

		if upperBound > 0 && count > upperBound {
			warnings = append(warnings, CardinalityWarning{
				JobName:      j.Name,
				Expression:   j.Expression,
				FiresPerWeek: count,
				Message: fmt.Sprintf(
					"fires %d times/week, exceeds upper bound of %d",
					count, upperBound,
				),
			})
		} else if lowerBound > 0 && count < lowerBound {
			warnings = append(warnings, CardinalityWarning{
				JobName:      j.Name,
				Expression:   j.Expression,
				FiresPerWeek: count,
				Message: fmt.Sprintf(
					"fires only %d times/week, below lower bound of %d",
					count, lowerBound,
				),
			})
		}
	}

	return warnings
}

// firesPerWeek returns the number of distinct (dow, hour, minute) combinations
// encoded in the parsed schedule, giving an upper-bound weekly fire count.
func firesPerWeek(s *parser.Schedule) int {
	return len(s.DayOfWeek) * len(s.Hours) * len(s.Minutes)
}
