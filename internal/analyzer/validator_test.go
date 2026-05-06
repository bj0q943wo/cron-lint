package analyzer

import (
	"testing"
)

func TestValidateJobs_NoWarnings(t *testing.T) {
	jobs := []Job{
		{Name: "backup", Schedule: mustParse("0 2 * * *")},
		{Name: "report", Schedule: mustParse("30 8 * * 1")},
	}
	warnings := ValidateJobs(jobs)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %d: %+v", len(warnings), warnings)
	}
}

func TestValidateJobs_EveryMinute(t *testing.T) {
	jobs := []Job{
		{Name: "poller", Schedule: mustParse("* * * * *")},
	}
	warnings := ValidateJobs(jobs)
	if len(warnings) == 0 {
		t.Fatal("expected at least one warning for every-minute schedule")
	}
	found := false
	for _, w := range warnings {
		if w.Job.Name == "poller" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected warning for job 'poller'")
	}
}

func TestValidateJobs_UnreachableDay31(t *testing.T) {
	// February (2) never has 31 days.
	jobs := []Job{
		{Name: "feb-end", Schedule: mustParse("0 0 31 2 *")},
	}
	warnings := ValidateJobs(jobs)
	if len(warnings) == 0 {
		t.Fatal("expected a warning for day 31 in February")
	}
	for _, w := range warnings {
		if w.Job.Name == "feb-end" {
			return
		}
	}
	t.Error("warning not attributed to job 'feb-end'")
}

func TestValidateJobs_Day31InLongMonth(t *testing.T) {
	// January (1) has 31 days — should not warn.
	jobs := []Job{
		{Name: "jan-end", Schedule: mustParse("0 0 31 1 *")},
	}
	warnings := ValidateJobs(jobs)
	for _, w := range warnings {
		if w.Job.Name == "jan-end" {
			t.Errorf("unexpected warning for jan-end: %s", w.Message)
		}
	}
}

func TestValidateJobs_MultipleWarnings(t *testing.T) {
	jobs := []Job{
		{Name: "heavy", Schedule: mustParse("* * * * *")},
		{Name: "feb31", Schedule: mustParse("0 0 31 2 *")},
		{Name: "clean", Schedule: mustParse("0 3 * * 0")},
	}
	warnings := ValidateJobs(jobs)
	if len(warnings) < 2 {
		t.Errorf("expected at least 2 warnings, got %d", len(warnings))
	}
}
