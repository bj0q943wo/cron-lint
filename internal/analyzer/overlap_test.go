package analyzer_test

import (
	"testing"

	"github.com/cron-lint/internal/analyzer"
	"github.com/cron-lint/internal/parser"
)

func mustParse(t *testing.T, expr string) *parser.Schedule {
	t.Helper()
	s, err := parser.Parse(expr)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", expr, err)
	}
	return s
}

func TestDetectOverlaps_NoOverlap(t *testing.T) {
	jobs := []analyzer.Job{
		{Name: "backup", Expression: "0 2 * * *", Schedule: mustParse(t, "0 2 * * *")},
		{Name: "report", Expression: "0 3 * * *", Schedule: mustParse(t, "0 3 * * *")},
	}
	warnings := analyzer.DetectOverlaps(jobs)
	if len(warnings) != 0 {
		t.Errorf("expected no overlaps, got %d", len(warnings))
	}
}

func TestDetectOverlaps_ExactOverlap(t *testing.T) {
	jobs := []analyzer.Job{
		{Name: "jobA", Expression: "30 6 * * *", Schedule: mustParse(t, "30 6 * * *")},
		{Name: "jobB", Expression: "30 6 * * *", Schedule: mustParse(t, "30 6 * * *")},
	}
	warnings := analyzer.DetectOverlaps(jobs)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 overlap, got %d", len(warnings))
	}
	if warnings[0].JobA != "jobA" || warnings[0].JobB != "jobB" {
		t.Errorf("unexpected job names: %v", warnings[0])
	}
}

func TestDetectOverlaps_PartialMinuteOverlap(t *testing.T) {
	jobs := []analyzer.Job{
		{Name: "every5", Expression: "*/5 * * * *", Schedule: mustParse(t, "*/5 * * * *")},
		{Name: "every15", Expression: "*/15 * * * *", Schedule: mustParse(t, "*/15 * * * *")},
	}
	warnings := analyzer.DetectOverlaps(jobs)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 overlap, got %d", len(warnings))
	}
	// 0, 15, 30, 45 are common
	if len(warnings[0].Minutes) != 4 {
		t.Errorf("expected 4 overlapping minutes, got %v", warnings[0].Minutes)
	}
}

func TestDetectOverlaps_MultipleJobs(t *testing.T) {
	jobs := []analyzer.Job{
		{Name: "a", Expression: "0 * * * *", Schedule: mustParse(t, "0 * * * *")},
		{Name: "b", Expression: "0 * * * *", Schedule: mustParse(t, "0 * * * *")},
		{Name: "c", Expression: "0 * * * *", Schedule: mustParse(t, "0 * * * *")},
	}
	warnings := analyzer.DetectOverlaps(jobs)
	// pairs: (a,b), (a,c), (b,c) => 3
	if len(warnings) != 3 {
		t.Errorf("expected 3 overlaps, got %d", len(warnings))
	}
}
