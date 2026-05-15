package analyzer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func makeIdleJob(expr string) Job {
	j, err := mustParseSchedule(expr)
	if err != nil {
		panic(err)
	}
	return j
}

func TestCheckIdle_NoGap_AllHoursCovered(t *testing.T) {
	// "* * * * *" fires every minute of every hour — no gap possible.
	jobs := []Job{makeIdleJob("* * * * *")}
	got := CheckIdle(jobs, DefaultIdleOptions())
	if len(got) != 0 {
		t.Fatalf("expected no warnings, got %v", got)
	}
}

func TestCheckIdle_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{{Name: "bad", Schedule: nil}}
	got := CheckIdle(jobs, DefaultIdleOptions())
	if len(got) != 0 {
		t.Fatalf("expected no warnings for nil schedule, got %v", got)
	}
}

func TestCheckIdle_SingleJobNightGap(t *testing.T) {
	// Runs only at 09:00 — leaves a large overnight gap.
	jobs := []Job{makeIdleJob("0 9 * * *")}
	got := CheckIdle(jobs, IdleOptions{MinGapHours: 4})
	if len(got) == 0 {
		t.Fatal("expected at least one idle warning")
	}
	// The gap from hour 10 through 8 (next day) is 23 hours.
	var maxGap int
	for _, w := range got {
		if w.GapHours > maxGap {
			maxGap = w.GapHours
		}
	}
	if maxGap < 4 {
		t.Errorf("expected gap >= 4, got %d", maxGap)
	}
}

func TestCheckIdle_ThresholdRespected(t *testing.T) {
	// Two jobs 3 hours apart — should not warn at threshold=4.
	jobs := []Job{
		makeIdleJob("0 0 * * *"),
		makeIdleJob("0 3 * * *"),
	}
	got := CheckIdle(jobs, IdleOptions{MinGapHours: 4})
	// gap between hour 1 and 2 is only 2 hours; between 4 and 23 is 20 hours.
	for _, w := range got {
		if w.GapHours < 4 {
			t.Errorf("warning emitted for gap smaller than threshold: %v", w)
		}
	}
}

func TestCheckIdle_WarningFields(t *testing.T) {
	jobs := []Job{makeIdleJob("0 0 * * *")}
	got := CheckIdle(jobs, IdleOptions{MinGapHours: 4})
	if len(got) == 0 {
		t.Fatal("expected warnings")
	}
	w := got[0]
	if w.Message == "" {
		t.Error("expected non-empty message")
	}
	if w.GapHours <= 0 {
		t.Errorf("expected positive GapHours, got %d", w.GapHours)
	}
	_ = cmp.Diff(nil, nil) // ensure import used
}

func TestCheckIdle_DefaultOptions_ZeroMinGap(t *testing.T) {
	// Passing zero MinGapHours should fall back to default (4).
	jobs := []Job{makeIdleJob("0 0 * * *")}
	got := CheckIdle(jobs, IdleOptions{MinGapHours: 0})
	for _, w := range got {
		if w.GapHours < 4 {
			t.Errorf("gap %d is below default threshold 4", w.GapHours)
		}
	}
}
