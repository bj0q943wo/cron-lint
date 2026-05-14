package analyzer

import (
	"testing"

	"github.com/your-org/cron-lint/internal/parser"
)

func makeCardJob(name, expr string) Job {
	sched, err := parser.Parse(expr)
	if err != nil {
		panic("makeCardJob: " + err.Error())
	}
	return Job{Name: name, Expression: expr, Schedule: sched}
}

func TestCheckCardinality_WithinBounds(t *testing.T) {
	// "0 9 * * 1-5" fires 1 min × 1 hour × 5 days = 5 times/week
	jobs := []Job{makeCardJob("weekday-morning", "0 9 * * 1-5")}
	warnings := CheckCardinality(jobs, 1, 100)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

func TestCheckCardinality_ExceedsUpperBound(t *testing.T) {
	// "* * * * *" fires 60 × 24 × 7 = 10 080 times/week
	jobs := []Job{makeCardJob("every-minute", "* * * * *")}
	warnings := CheckCardinality(jobs, 0, 500)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].JobName != "every-minute" {
		t.Errorf("unexpected job name: %s", warnings[0].JobName)
	}
	if warnings[0].FiresPerWeek != 60*24*7 {
		t.Errorf("unexpected fires/week: %d", warnings[0].FiresPerWeek)
	}
}

func TestCheckCardinality_BelowLowerBound(t *testing.T) {
	// "0 0 * * 0" fires 1 time/week (Sunday midnight)
	jobs := []Job{makeCardJob("weekly-sunday", "0 0 * * 0")}
	warnings := CheckCardinality(jobs, 5, 0)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].FiresPerWeek != 1 {
		t.Errorf("unexpected fires/week: %d", warnings[0].FiresPerWeek)
	}
}

func TestCheckCardinality_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{{Name: "broken", Expression: "bad", Schedule: nil}}
	warnings := CheckCardinality(jobs, 1, 100)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings for nil schedule, got %d", len(warnings))
	}
}

func TestCheckCardinality_MultipleJobs(t *testing.T) {
	jobs := []Job{
		makeCardJob("ok", "0 9 * * 1-5"),        // 5/week — within [1, 100]
		makeCardJob("too-frequent", "* * * * *"), // 10080/week — exceeds 100
		makeCardJob("too-rare", "0 0 * * 0"),     // 1/week — below 3
	}
	warnings := CheckCardinality(jobs, 3, 100)
	if len(warnings) != 2 {
		t.Fatalf("expected 2 warnings, got %d", len(warnings))
	}
	names := map[string]bool{}
	for _, w := range warnings {
		names[w.JobName] = true
	}
	if !names["too-frequent"] || !names["too-rare"] {
		t.Errorf("wrong jobs warned: %v", names)
	}
}
