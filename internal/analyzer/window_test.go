package analyzer

import (
	"testing"

	"github.com/your-org/cron-lint/internal/parser"
)

var maintenanceWindows = []TimeWindow{
	{Name: "nightly", Start: 2 * 60, End: 4 * 60}, // 02:00–04:00
}

func makeWindowJob(name, expr string) parser.Job {
	j, err := parser.Parse(expr)
	if err != nil {
		panic(err)
	}
	return parser.Job{Name: name, Expression: expr, Schedule: j}
}

func TestCheckWindows_NoConflict(t *testing.T) {
	jobs := []parser.Job{
		makeWindowJob("morning", "30 8 * * *"),
	}
	got := CheckWindows(jobs, maintenanceWindows)
	if len(got) != 0 {
		t.Fatalf("expected 0 warnings, got %d", len(got))
	}
}

func TestCheckWindows_InsideWindow(t *testing.T) {
	jobs := []parser.Job{
		makeWindowJob("backup", "15 3 * * *"),
	}
	got := CheckWindows(jobs, maintenanceWindows)
	if len(got) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(got))
	}
	if got[0].Job.Name != "backup" {
		t.Errorf("unexpected job name %q", got[0].Job.Name)
	}
	if got[0].Window.Name != "nightly" {
		t.Errorf("unexpected window %q", got[0].Window.Name)
	}
}

func TestCheckWindows_PartiallyOutside(t *testing.T) {
	// Runs every hour — some hours outside the window.
	jobs := []parser.Job{
		makeWindowJob("hourly", "0 * * * *"),
	}
	got := CheckWindows(jobs, maintenanceWindows)
	if len(got) != 0 {
		t.Fatalf("expected 0 warnings, got %d", len(got))
	}
}

func TestCheckWindows_MultipleJobs(t *testing.T) {
	jobs := []parser.Job{
		makeWindowJob("safe", "0 9 * * *"),
		makeWindowJob("risky", "0 2 * * *"),
		makeWindowJob("also-risky", "45 3 * * *"),
	}
	got := CheckWindows(jobs, maintenanceWindows)
	if len(got) != 2 {
		t.Fatalf("expected 2 warnings, got %d", len(got))
	}
}

func TestCheckWindows_NilScheduleSkipped(t *testing.T) {
	jobs := []parser.Job{
		{Name: "broken", Expression: "bad", Schedule: nil},
	}
	got := CheckWindows(jobs, maintenanceWindows)
	if len(got) != 0 {
		t.Fatalf("expected 0 warnings for nil schedule, got %d", len(got))
	}
}

func TestFormatWindowWarnings_Empty(t *testing.T) {
	out := FormatWindowWarnings(nil)
	if out != "no window conflicts detected" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatWindowWarnings_WithWarnings(t *testing.T) {
	warnings := []WindowWarning{
		{Job: parser.Job{Name: "backup"}, Window: maintenanceWindows[0], Message: "job \"backup\" always fires inside window \"nightly\" (02:00–04:00)"},
	}
	out := FormatWindowWarnings(warnings)
	if out == "" || out == "no window conflicts detected" {
		t.Errorf("expected non-empty warning output, got: %q", out)
	}
}
