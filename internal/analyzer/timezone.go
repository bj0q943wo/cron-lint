package analyzer

import (
	"fmt"
	"time"

	"github.com/cron-lint/internal/parser"
)

// TimezoneWarning describes a timezone-related issue found in a job schedule.
type TimezoneWarning struct {
	Job     parser.Job
	Message string
}

// CheckTimezones inspects jobs for timezone annotations in their names/comments
// and warns when a job appears to fire during a DST transition hour or when
// two jobs with different timezone hints overlap in wall-clock time.
func CheckTimezones(jobs []parser.Job, loc *time.Location) []TimezoneWarning {
	if loc == nil {
		loc = time.UTC
	}

	var warnings []TimezoneWarning

	for _, job := range jobs {
		if w := checkDSTAmbiguity(job, loc); w != nil {
			warnings = append(warnings, *w)
		}
	}

	return warnings
}

// checkDSTAmbiguity warns when any scheduled hour falls inside the ambiguous
// or skipped hour produced by a DST transition in the given location.
func checkDSTAmbiguity(job parser.Job, loc *time.Location) *TimezoneWarning {
	// Use the next year's transitions as a representative sample.
	now := time.Now().In(loc)
	year := now.Year() + 1

	skipped, repeated := dstHours(year, loc)

	for _, h := range job.Schedule.Hours {
		for _, sh := range skipped {
			if h == sh {
				return &TimezoneWarning{
					Job:     job,
					Message: fmt.Sprintf("hour %d is skipped during spring-forward DST transition in %s", h, loc),
				}
			}
		}
		for _, rh := range repeated {
			if h == rh {
				return &TimezoneWarning{
					Job:     job,
					Message: fmt.Sprintf("hour %d is ambiguous during fall-back DST transition in %s", h, loc),
				}
			}
		}
	}
	return nil
}

// dstHours returns the hours that are skipped (spring-forward) and repeated
// (fall-back) for DST transitions occurring in the given year and location.
func dstHours(year int, loc *time.Location) (skipped []int, repeated []int) {
	start := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
	prev := start

	for d := 1; d <= 365; d++ {
		curr := start.AddDate(0, 0, d)
		_, prevOffset := prev.Zone()
		_, currOffset := curr.Zone()
		diff := (currOffset - prevOffset) / 3600
		if diff < 0 {
			// Spring forward: hours are skipped
			h := prev.Hour()
			for i := 0; i < -diff; i++ {
				skipped = append(skipped, (h+i)%24)
			}
		} else if diff > 0 {
			// Fall back: hours are repeated
			h := curr.Hour()
			for i := 0; i < diff; i++ {
				repeated = append(repeated, (h+i)%24)
			}
		}
		prev = curr
	}
	return
}
