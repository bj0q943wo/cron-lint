package analyzer

import (
	"testing"

	"github.com/your-org/cron-lint/internal/parser"
)

func makeDriftJob(name, expr string) Job {
	sched, err := parser.Parse(expr)
	if err != nil {
		panic("makeDriftJob: " + err.Error())
	}
	return Job{Name: name, Expression: expr, Schedule: sched}
}

func TestCheckDrift_NoWarnings_SameMinutes(t *testing.T) {
	jobs := []Job{
		makeDriftJob("a", "0 9 * * *"),
		makeDriftJob("b", "0 9 * * *"),
	}
	warnings := CheckDrift(jobs, 5)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

func TestCheckDrift_NoWarnings_DifferentHours(t *testing.T) {
	jobs := []Job{
		makeDriftJob("a", "5 8 * * *"),
		makeDriftJob("b", "5 9 * * *"),
	}
	warnings := CheckDrift(jobs, 5)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings for different hours, got %d", len(warnings))
	}
}

func TestCheckDrift_OffsetExceedsThreshold(t *testing.T) {
	jobs := []Job{
		makeDriftJob("sync-db", "0 2 * * *"),
		makeDriftJob("sync-cache", "15 2 * * *"),
	}
	warnings := CheckDrift(jobs, 5)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	w := warnings[0]
	if w.OffsetMin != 15 {
		t.Errorf("expected offset 15, got %d", w.OffsetMin)
	}
	if w.JobA != "sync-db" || w.JobB != "sync-cache" {
		t.Errorf("unexpected job names: %q %q", w.JobA, w.JobB)
	}
}

func TestCheckDrift_OffsetWithinThreshold(t *testing.T) {
	jobs := []Job{
		makeDriftJob("a", "0 3 * * *"),
		makeDriftJob("b", "3 3 * * *"),
	}
	warnings := CheckDrift(jobs, 5)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings within threshold, got %d", len(warnings))
	}
}

func TestCheckDrift_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{
		{Name: "bad", Expression: "invalid", Schedule: nil},
		makeDriftJob("good", "0 6 * * *"),
	}
	warnings := CheckDrift(jobs, 5)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings when one schedule is nil, got %d", len(warnings))
	}
}

func TestCheckDrift_MultipleOffendingPairs(t *testing.T) {
	jobs := []Job{
		makeDriftJob("x", "0 12 * * *"),
		makeDriftJob("y", "20 12 * * *"),
		makeDriftJob("z", "40 12 * * *"),
	}
	warnings := CheckDrift(jobs, 5)
	if len(warnings) != 3 {
		t.Fatalf("expected 3 warnings, got %d", len(warnings))
	}
}
