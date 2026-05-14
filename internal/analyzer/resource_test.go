package analyzer

import (
	"testing"

	"github.com/your-org/cron-lint/internal/parser"
)

func makeResourceJob(name, expr string) Job {
	sched, err := parser.Parse(expr)
	if err != nil {
		panic("makeResourceJob: " + err.Error())
	}
	return Job{Name: name, Raw: expr, Schedule: sched}
}

func TestCheckResourceContention_NoContention(t *testing.T) {
	jobs := []Job{
		makeResourceJob("alpha", "5 8 * * *"),
		makeResourceJob("beta", "10 9 * * *"),
		makeResourceJob("gamma", "15 10 * * *"),
	}
	warns := CheckResourceContention(jobs, ResourceParams{MaxJobsPerSlot: 3})
	if len(warns) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warns))
	}
}

func TestCheckResourceContention_ExactThreshold(t *testing.T) {
	// Three jobs at exactly the threshold — no warning expected.
	jobs := []Job{
		makeResourceJob("a", "0 6 * * *"),
		makeResourceJob("b", "0 6 * * *"),
		makeResourceJob("c", "0 6 * * *"),
	}
	warns := CheckResourceContention(jobs, ResourceParams{MaxJobsPerSlot: 3})
	if len(warns) != 0 {
		t.Fatalf("expected no warnings at threshold, got %d", len(warns))
	}
}

func TestCheckResourceContention_ExceedsThreshold(t *testing.T) {
	// Four jobs at the same slot — all four should get a warning.
	jobs := []Job{
		makeResourceJob("j1", "0 6 * * *"),
		makeResourceJob("j2", "0 6 * * *"),
		makeResourceJob("j3", "0 6 * * *"),
		makeResourceJob("j4", "0 6 * * *"),
	}
	warns := CheckResourceContention(jobs, ResourceParams{MaxJobsPerSlot: 3})
	if len(warns) != 4 {
		t.Fatalf("expected 4 warnings (one per job), got %d", len(warns))
	}
	for _, w := range warns {
		if w.Slot != "06:00" {
			t.Errorf("unexpected slot %q", w.Slot)
		}
		if len(w.PeakJobs) != 4 {
			t.Errorf("expected 4 peak jobs, got %d", len(w.PeakJobs))
		}
	}
}

func TestCheckResourceContention_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{
		{Name: "bad", Raw: "invalid", Schedule: nil},
		makeResourceJob("ok", "0 6 * * *"),
	}
	warns := CheckResourceContention(jobs, ResourceParams{MaxJobsPerSlot: 1})
	// Only one valid job at 06:00 — no contention.
	if len(warns) != 0 {
		t.Fatalf("nil schedule should be skipped, got %d warnings", len(warns))
	}
}

func TestCheckResourceContention_DefaultParams(t *testing.T) {
	// Zero-value params should default MaxJobsPerSlot to 3.
	jobs := []Job{
		makeResourceJob("x1", "30 12 * * *"),
		makeResourceJob("x2", "30 12 * * *"),
		makeResourceJob("x3", "30 12 * * *"),
		makeResourceJob("x4", "30 12 * * *"),
	}
	warns := CheckResourceContention(jobs, ResourceParams{})
	if len(warns) == 0 {
		t.Fatal("expected warnings with default threshold of 3")
	}
}
