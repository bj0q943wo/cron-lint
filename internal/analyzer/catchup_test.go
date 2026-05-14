package analyzer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func makeCatchupJob(name, expr string) Job {
	return makeRetryJob(name, expr)
}

func TestCheckCatchup_FrequentJob_Warning(t *testing.T) {
	jobs := []Job{
		makeCatchupJob("every-minute", "* * * * *"),
	}
	got := CheckCatchup(jobs, 8, 10)
	if len(got) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(got))
	}
	if got[0].Job.Name != "every-minute" {
		t.Errorf("unexpected job name %q", got[0].Job.Name)
	}
	if got[0].EstimatedCatchup < 10 {
		t.Errorf("expected estimated catchup >= 10, got %d", got[0].EstimatedCatchup)
	}
}

func TestCheckCatchup_RareJob_NoWarning(t *testing.T) {
	// Runs once a week — far below any reasonable threshold.
	jobs := []Job{
		makeCatchupJob("weekly", "0 9 * * 1"),
	}
	got := CheckCatchup(jobs, 8, 10)
	if len(got) != 0 {
		t.Errorf("expected no warnings, got %d", len(got))
	}
}

func TestCheckCatchup_NilSchedule_Skipped(t *testing.T) {
	jobs := []Job{{Name: "nil-sched", Expression: "bad", Schedule: nil}}
	got := CheckCatchup(jobs, 8, 10)
	if len(got) != 0 {
		t.Errorf("expected no warnings for nil schedule, got %d", len(got))
	}
}

func TestCheckCatchup_DefaultParameters(t *testing.T) {
	// Passing 0 should fall back to defaults without panicking.
	jobs := []Job{
		makeCatchupJob("every-minute", "* * * * *"),
	}
	got := CheckCatchup(jobs, 0, 0)
	if len(got) == 0 {
		t.Error("expected at least one warning with default parameters")
	}
}

func TestCheckCatchup_HourlyJob_BelowThreshold(t *testing.T) {
	// 24 fires/day → 8 h outage → ~8 missed — below default threshold of 10.
	jobs := []Job{
		makeCatchupJob("hourly", "0 * * * *"),
	}
	got := CheckCatchup(jobs, 8, 10)
	if len(got) != 0 {
		t.Errorf("expected no warnings for hourly job, got %d", len(got))
	}
}

func TestCheckCatchup_MultipleJobs_OnlyFrequentWarned(t *testing.T) {
	jobs := []Job{
		makeCatchupJob("every-minute", "* * * * *"),
		makeCatchupJob("daily", "0 6 * * *"),
	}
	got := CheckCatchup(jobs, 8, 10)
	if len(got) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(got))
	}
	if diff := cmp.Diff("every-minute", got[0].Job.Name); diff != "" {
		t.Errorf("job name mismatch (-want +got):\n%s", diff)
	}
}
