package analyzer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func makeBurstJob(name, expr string) JobEntry {
	return mustParse(name, expr)
}

func TestCheckBursts_NoBurst(t *testing.T) {
	jobs := []JobEntry{
		makeBurstJob("a", "0 * * * *"),
		makeBurstJob("b", "30 * * * *"),
		makeBurstJob("c", "15 * * * *"),
	}
	got := CheckBursts(jobs, 5, 3)
	if len(got) != 0 {
		t.Errorf("expected no bursts, got %v", got)
	}
}

func TestCheckBursts_ExactThreshold(t *testing.T) {
	jobs := []JobEntry{
		makeBurstJob("a", "0 * * * *"),
		makeBurstJob("b", "1 * * * *"),
		makeBurstJob("c", "2 * * * *"),
	}
	got := CheckBursts(jobs, 5, 3)
	if len(got) != 1 {
		t.Fatalf("expected 1 burst warning, got %d", len(got))
	}
	if got[0].Count != 3 {
		t.Errorf("expected count 3, got %d", got[0].Count)
	}
	if got[0].Minute != 0 {
		t.Errorf("expected start minute 0, got %d", got[0].Minute)
	}
}

func TestCheckBursts_BelowThreshold(t *testing.T) {
	jobs := []JobEntry{
		makeBurstJob("a", "0 * * * *"),
		makeBurstJob("b", "1 * * * *"),
	}
	got := CheckBursts(jobs, 5, 3)
	if len(got) != 0 {
		t.Errorf("expected no bursts, got %v", got)
	}
}

func TestCheckBursts_NilScheduleSkipped(t *testing.T) {
	jobs := []JobEntry{
		{Name: "nil-sched", Expression: "bad", Schedule: nil},
		makeBurstJob("a", "5 * * * *"),
		makeBurstJob("b", "6 * * * *"),
		makeBurstJob("c", "7 * * * *"),
	}
	got := CheckBursts(jobs, 5, 3)
	if len(got) != 1 {
		t.Fatalf("expected 1 burst, got %d", len(got))
	}
	for _, j := range got[0].Jobs {
		if j.Name == "nil-sched" {
			t.Error("nil-schedule job should not appear in burst group")
		}
	}
}

func TestCheckBursts_WindowsNotDoubleReported(t *testing.T) {
	// Jobs at minutes 0,1,2,3 — only one warning should be emitted.
	jobs := []JobEntry{
		makeBurstJob("a", "0 * * * *"),
		makeBurstJob("b", "1 * * * *"),
		makeBurstJob("c", "2 * * * *"),
		makeBurstJob("d", "3 * * * *"),
	}
	got := CheckBursts(jobs, 5, 3)
	if len(got) != 1 {
		t.Errorf("expected 1 merged warning, got %d: %v", len(got), got)
	}
	_ = cmp.Diff // ensure import used
}

func TestCheckBursts_InvalidParams(t *testing.T) {
	jobs := []JobEntry{makeBurstJob("a", "0 * * * *")}
	if got := CheckBursts(jobs, 0, 3); got != nil {
		t.Error("zero windowMinutes should return nil")
	}
	if got := CheckBursts(jobs, 5, 0); got != nil {
		t.Error("zero threshold should return nil")
	}
}
