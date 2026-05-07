package analyzer

import (
	"testing"
	"time"

	"github.com/user/cron-lint/internal/parser"
)

// reference is a fixed Monday 2024-03-04 00:00 UTC used across all tests.
var staleRef = time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)

func makeStaleJob(name, expr string) Job {
	sched, err := parser.Parse(expr)
	if err != nil {
		panic("makeStaleJob: " + err.Error())
	}
	return Job{Name: name, Expression: expr, Schedule: sched}
}

func TestCheckStaleness_NormalJob_NoWarning(t *testing.T) {
	jobs := []Job{makeStaleJob("daily", "0 9 * * *")}
	warnings := CheckStaleness(jobs, staleRef)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

func TestCheckStaleness_YearlyJob_NoWarning(t *testing.T) {
	// Fires on Jan 1 — within 366-day window.
	jobs := []Job{makeStaleJob("yearly", "0 9 1 1 *")}
	warnings := CheckStaleness(jobs, staleRef)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings for yearly job, got %d", len(warnings))
	}
}

func TestCheckStaleness_ImpossibleDay_Warning(t *testing.T) {
	// Feb 30 never exists.
	jobs := []Job{
		{
			Name:       "impossible",
			Expression: "0 0 30 2 *",
			Schedule: &parser.Schedule{
				Minutes:     []int{0},
				Hours:       []int{0},
				DaysOfMonth: []int{30},
				Months:      []int{2},
				DaysOfWeek:  []int{0, 1, 2, 3, 4, 5, 6},
			},
		},
	}
	warnings := CheckStaleness(jobs, staleRef)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].JobName != "impossible" {
		t.Errorf("unexpected job name %q", warnings[0].JobName)
	}
}

func TestCheckStaleness_NilSchedule_Skipped(t *testing.T) {
	jobs := []Job{{Name: "broken", Expression: "", Schedule: nil}}
	warnings := CheckStaleness(jobs, staleRef)
	if len(warnings) != 0 {
		t.Fatalf("nil schedule should be skipped, got %d warnings", len(warnings))
	}
}

func TestCheckStaleness_MultipleJobs_OnlyStaleWarned(t *testing.T) {
	goodJob := makeStaleJob("good", "*/5 * * * *")
	badJob := Job{
		Name:       "stale",
		Expression: "0 0 31 2 *",
		Schedule: &parser.Schedule{
			Minutes:     []int{0},
			Hours:       []int{0},
			DaysOfMonth: []int{31},
			Months:      []int{2},
			DaysOfWeek:  []int{0, 1, 2, 3, 4, 5, 6},
		},
	}
	warnings := CheckStaleness(jobs(goodJob, badJob), staleRef)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].JobName != "stale" {
		t.Errorf("expected stale job warned, got %q", warnings[0].JobName)
	}
}

func jobs(jj ...Job) []Job { return jj }
