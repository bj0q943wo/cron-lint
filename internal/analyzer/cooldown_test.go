package analyzer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func makeCooldownJob(name, expr string) Job {
	sched, err := mustParseSchedule(expr)
	if err != nil {
		panic(err)
	}
	return Job{Name: name, Expression: expr, Schedule: sched}
}

func TestCheckCooldown_NoWarning_LargeGap(t *testing.T) {
	jobs := []Job{
		makeCooldownJob("a", "0 * * * *"),  // fires at :00
		makeCooldownJob("b", "30 * * * *"), // fires at :30
	}
	got := CheckCooldown(jobs, DefaultCooldownOptions())
	if len(got) != 0 {
		t.Errorf("expected no warnings, got %v", got)
	}
}

func TestCheckCooldown_ExactThreshold_NoWarning(t *testing.T) {
	opts := CooldownOptions{MinGapMinutes: 5}
	jobs := []Job{
		makeCooldownJob("a", "0 * * * *"),
		makeCooldownJob("b", "5 * * * *"), // exactly 5 minutes apart — not a warning
	}
	got := CheckCooldown(jobs, opts)
	if len(got) != 0 {
		t.Errorf("expected no warnings at exact threshold, got %v", got)
	}
}

func TestCheckCooldown_BelowThreshold_Warning(t *testing.T) {
	opts := CooldownOptions{MinGapMinutes: 5}
	jobs := []Job{
		makeCooldownJob("a", "0 * * * *"),
		makeCooldownJob("b", "3 * * * *"), // only 3 minutes apart
	}
	got := CheckCooldown(jobs, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(got))
	}
	if got[0].GapMinutes != 3 {
		t.Errorf("expected gap=3, got %d", got[0].GapMinutes)
	}
	if got[0].JobA != "a" || got[0].JobB != "b" {
		t.Errorf("unexpected job names: %q %q", got[0].JobA, got[0].JobB)
	}
}

func TestCheckCooldown_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{
		{Name: "nil-sched", Expression: "* * * * *", Schedule: nil},
		makeCooldownJob("b", "1 * * * *"),
	}
	got := CheckCooldown(jobs, DefaultCooldownOptions())
	if len(got) != 0 {
		t.Errorf("expected no warnings when one schedule is nil, got %v", got)
	}
}

func TestCheckCooldown_MultipleViolations(t *testing.T) {
	opts := CooldownOptions{MinGapMinutes: 10}
	jobs := []Job{
		makeCooldownJob("a", "0 * * * *"),
		makeCooldownJob("b", "2 * * * *"),  // 2 min from a
		makeCooldownJob("c", "15 * * * *"), // 15 min from a — ok
		makeCooldownJob("d", "4 * * * *"),  // 2 min from b, 4 min from a
	}
	got := CheckCooldown(jobs, opts)
	// Expect violations: (a,b), (a,d), (b,d)
	if len(got) < 2 {
		t.Errorf("expected at least 2 warnings, got %d: %v", len(got), got)
	}
	_ = cmp.Diff // ensure import used
}

func TestCheckCooldown_DefaultOptions_Applied(t *testing.T) {
	// MinGapMinutes=0 should fall back to default (5)
	opts := CooldownOptions{MinGapMinutes: 0}
	jobs := []Job{
		makeCooldownJob("a", "0 * * * *"),
		makeCooldownJob("b", "3 * * * *"),
	}
	got := CheckCooldown(jobs, opts)
	if len(got) != 1 {
		t.Errorf("expected default min-gap to trigger warning, got %d warnings", len(got))
	}
}
