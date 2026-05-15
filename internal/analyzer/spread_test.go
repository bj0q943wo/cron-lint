package analyzer

import (
	"testing"

	"github.com/your-org/cron-lint/internal/parser"
)

func makeSpreadJob(name, expr string) Job {
	sched, err := parser.Parse(expr)
	if err != nil {
		panic(err)
	}
	return Job{Name: name, Expression: expr, Schedule: sched}
}

func TestCheckSpread_NoCluster(t *testing.T) {
	jobs := []Job{
		makeSpreadJob("a", "0 * * * *"),
		makeSpreadJob("b", "15 * * * *"),
		makeSpreadJob("c", "30 * * * *"),
		makeSpreadJob("d", "45 * * * *"),
	}
	warn := CheckSpread(jobs, DefaultSpreadOptions())
	if len(warn) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warn))
	}
}

func TestCheckSpread_ExactThreshold(t *testing.T) {
	// 3 jobs within 5 minutes — exactly at the limit, no warning expected.
	jobs := []Job{
		makeSpreadJob("a", "0 * * * *"),
		makeSpreadJob("b", "1 * * * *"),
		makeSpreadJob("c", "2 * * * *"),
	}
	warn := CheckSpread(jobs, DefaultSpreadOptions())
	if len(warn) != 0 {
		t.Fatalf("expected no warnings at threshold, got %d", len(warn))
	}
}

func TestCheckSpread_ExceedsThreshold(t *testing.T) {
	// 4 jobs within 5 minutes — exceeds MaxJobsInWindow=3.
	jobs := []Job{
		makeSpreadJob("a", "0 * * * *"),
		makeSpreadJob("b", "1 * * * *"),
		makeSpreadJob("c", "2 * * * *"),
		makeSpreadJob("d", "3 * * * *"),
	}
	warn := CheckSpread(jobs, DefaultSpreadOptions())
	if len(warn) == 0 {
		t.Fatal("expected at least one warning")
	}
	if len(warn[0].Jobs) != 4 {
		t.Errorf("expected 4 jobs in warning, got %d", len(warn[0].Jobs))
	}
}

func TestCheckSpread_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{
		{Name: "nil", Expression: "bad", Schedule: nil},
		makeSpreadJob("a", "5 * * * *"),
	}
	warn := CheckSpread(jobs, DefaultSpreadOptions())
	if len(warn) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warn))
	}
}

func TestCheckSpread_CustomOptions(t *testing.T) {
	opts := SpreadOptions{WindowSize: 10, MaxJobsInWindow: 2}
	jobs := []Job{
		makeSpreadJob("a", "0 * * * *"),
		makeSpreadJob("b", "5 * * * *"),
		makeSpreadJob("c", "9 * * * *"),
	}
	warn := CheckSpread(jobs, opts)
	if len(warn) == 0 {
		t.Fatal("expected warning with custom options")
	}
}

func TestCheckSpread_WrapAround(t *testing.T) {
	// Jobs at minute 58, 59, 0, 1 should cluster when window wraps.
	jobs := []Job{
		makeSpreadJob("a", "58 * * * *"),
		makeSpreadJob("b", "59 * * * *"),
		makeSpreadJob("c", "0 * * * *"),
		makeSpreadJob("d", "1 * * * *"),
	}
	warn := CheckSpread(jobs, DefaultSpreadOptions())
	if len(warn) == 0 {
		t.Fatal("expected wrap-around cluster warning")
	}
}
