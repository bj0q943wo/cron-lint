package analyzer

import (
	"testing"

	"github.com/your-org/cron-lint/internal/parser"
)

func makeThrottleJob(name, expr string) Job {
	sched, err := parser.Parse(expr)
	if err != nil {
		panic(err)
	}
	return Job{Name: name, Expression: expr, Schedule: sched}
}

func TestCheckThrottle_NoBurst(t *testing.T) {
	jobs := []Job{
		makeThrottleJob("a", "0 * * * *"),
		makeThrottleJob("b", "30 * * * *"),
	}
	opts := ThrottleOptions{MaxFiringsPer5Min: 5}
	warnings := CheckThrottle(jobs, opts)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d: %v", len(warnings), warnings)
	}
}

func TestCheckThrottle_ExactThreshold(t *testing.T) {
	// 5 jobs all firing at minute 0 of every hour → 5 firings per slot, threshold 5 → no warning
	var jobs []Job
	for i := 0; i < 5; i++ {
		jobs = append(jobs, makeThrottleJob("j", "0 * * * *"))
	}
	opts := ThrottleOptions{MaxFiringsPer5Min: 5}
	warnings := CheckThrottle(jobs, opts)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings at threshold boundary, got %d", len(warnings))
	}
}

func TestCheckThrottle_ExceedsThreshold(t *testing.T) {
	// 6 distinct jobs all firing at minute 0 → 6 firings per (h,block) pair
	var jobs []Job
	for i := 0; i < 6; i++ {
		jobs = append(jobs, makeThrottleJob("job", "0 * * * *"))
	}
	opts := ThrottleOptions{MaxFiringsPer5Min: 5}
	warnings := CheckThrottle(jobs, opts)
	// Expect 24 warnings (one per hour)
	if len(warnings) != 24 {
		t.Fatalf("expected 24 warnings (one per hour), got %d", len(warnings))
	}
	for _, w := range warnings {
		if w.Firings != 6 {
			t.Errorf("expected 6 firings, got %d", w.Firings)
		}
		if w.Threshold != 5 {
			t.Errorf("expected threshold 5, got %d", w.Threshold)
		}
	}
}

func TestCheckThrottle_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{
		{Name: "bad", Expression: "* * * * *", Schedule: nil},
	}
	opts := DefaultThrottleOptions
	warnings := CheckThrottle(jobs, opts)
	if len(warnings) != 0 {
		t.Fatalf("nil schedule should be skipped, got %d warnings", len(warnings))
	}
}

func TestCheckThrottle_MultipleJobsDifferentSlots(t *testing.T) {
	// Jobs spread across different 5-min blocks – no single block overflows
	jobs := []Job{
		makeThrottleJob("a", "0 0 * * *"),
		makeThrottleJob("b", "10 0 * * *"),
		makeThrottleJob("c", "20 0 * * *"),
	}
	opts := ThrottleOptions{MaxFiringsPer5Min: 2}
	warnings := CheckThrottle(jobs, opts)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}
