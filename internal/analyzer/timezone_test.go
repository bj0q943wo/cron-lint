package analyzer

import (
	"testing"
	"time"
)

func TestCheckTimezones_UTC_NoWarnings(t *testing.T) {
	// UTC has no DST transitions, so no warnings should ever be produced.
	jobs := []mustParseJob{
		{expr: "0 2 * * *", name: "backup"},
		{expr: "30 10 * * *", name: "report"},
	}
	parsedJobs := mustParseJobs(t, jobs)

	warnings := CheckTimezones(parsedJobs, time.UTC)
	if len(warnings) != 0 {
		t.Errorf("expected 0 warnings in UTC, got %d: %v", len(warnings), warnings)
	}
}

func TestCheckTimezones_NilLocation_NoWarnings(t *testing.T) {
	jobs := mustParseJobs(t, []mustParseJob{{expr: "0 2 * * *", name: "job"}})
	warnings := CheckTimezones(jobs, nil)
	if len(warnings) != 0 {
		t.Errorf("expected 0 warnings with nil loc, got %d", len(warnings))
	}
}

func TestCheckTimezones_DSTLocation_AmbiguousHour(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("timezone data not available")
	}

	// America/New_York falls back at 2am, making hour 1 ambiguous.
	jobs := mustParseJobs(t, []mustParseJob{{expr: "30 1 * * *", name: "nightly"}})
	warnings := CheckTimezones(jobs, loc)

	if len(warnings) == 0 {
		t.Error("expected a DST ambiguity warning for hour 1 in America/New_York, got none")
	}
}

func TestCheckTimezones_DSTLocation_SkippedHour(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("timezone data not available")
	}

	// America/New_York springs forward at 2am, skipping hour 2.
	jobs := mustParseJobs(t, []mustParseJob{{expr: "0 2 * * *", name: "sync"}})
	warnings := CheckTimezones(jobs, loc)

	if len(warnings) == 0 {
		t.Error("expected a DST skipped-hour warning for hour 2 in America/New_York, got none")
	}
}

func TestCheckTimezones_MultipleJobs_OnlyAffectedWarned(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("timezone data not available")
	}

	jobs := mustParseJobs(t, []mustParseJob{
		{expr: "0 2 * * *", name: "dst-affected"},
		{expr: "0 12 * * *", name: "safe-noon"},
	})
	warnings := CheckTimezones(jobs, loc)

	for _, w := range warnings {
		if w.Job.Name == "safe-noon" {
			t.Errorf("unexpected warning for safe-noon job: %s", w.Message)
		}
	}
}

// mustParseJob is a helper struct for table-driven tests.
type mustParseJob struct {
	expr string
	name string
}

func mustParseJobs(t *testing.T, specs []mustParseJob) []mustParseJobResult {
	t.Helper()
	result := make([]mustParseJobResult, 0, len(specs))
	for _, s := range specs {
		j := mustParse(t, s.expr)
		j.Name = s.name
		result = append(result, j)
	}
	return result
}
