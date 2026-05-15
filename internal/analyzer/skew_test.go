package analyzer

import (
	"testing"

	"github.com/your-org/cron-lint/internal/parser"
)

func makeSkewJob(name, expr string) Job {
	s, err := parser.Parse(expr)
	if err != nil {
		panic(err)
	}
	return Job{Name: name, Expression: expr, Schedule: s}
}

func TestCheckSkew_NoConflict(t *testing.T) {
	jobs := []Job{
		makeSkewJob("a", "0 9 * * 1"),
		makeSkewJob("b", "30 10 * * 2"),
	}
	ws := CheckSkew(jobs)
	if len(ws) != 0 {
		t.Fatalf("expected no warnings, got %d", len(ws))
	}
}

func TestCheckSkew_TwoJobsSameClock(t *testing.T) {
	jobs := []Job{
		makeSkewJob("alpha", "0 9 * * 1"),
		makeSkewJob("beta", "0 9 * * 3"),
	}
	ws := CheckSkew(jobs)
	if len(ws) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(ws))
	}
	if ws[0].Hour != 9 || ws[0].Minute != 0 {
		t.Errorf("unexpected clock time %02d:%02d", ws[0].Hour, ws[0].Minute)
	}
	if len(ws[0].Jobs) != 2 {
		t.Errorf("expected 2 job names, got %v", ws[0].Jobs)
	}
}

func TestCheckSkew_ThreeJobsSameClock(t *testing.T) {
	jobs := []Job{
		makeSkewJob("j1", "15 6 * * *"),
		makeSkewJob("j2", "15 6 * * *"),
		makeSkewJob("j3", "15 6 * * *"),
	}
	ws := CheckSkew(jobs)
	if len(ws) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(ws))
	}
	if len(ws[0].Jobs) != 3 {
		t.Errorf("expected 3 jobs, got %v", ws[0].Jobs)
	}
}

func TestCheckSkew_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{
		{Name: "nil-sched", Expression: "bad", Schedule: nil},
		makeSkewJob("ok", "0 8 * * *"),
	}
	ws := CheckSkew(jobs)
	if len(ws) != 0 {
		t.Fatalf("expected no warnings, got %d", len(ws))
	}
}

func TestCheckSkew_MultipleClockTimes(t *testing.T) {
	jobs := []Job{
		makeSkewJob("x1", "0 8 * * *"),
		makeSkewJob("x2", "0 8 * * *"),
		makeSkewJob("y1", "30 12 * * *"),
		makeSkewJob("y2", "30 12 * * *"),
	}
	ws := CheckSkew(jobs)
	if len(ws) != 2 {
		t.Fatalf("expected 2 warnings, got %d", len(ws))
	}
	if ws[0].Hour != 8 {
		t.Errorf("expected first warning at hour 8, got %d", ws[0].Hour)
	}
	if ws[1].Hour != 12 {
		t.Errorf("expected second warning at hour 12, got %d", ws[1].Hour)
	}
}

func TestFormatSkewWarnings_Empty(t *testing.T) {
	lines := FormatSkewWarnings(nil)
	if len(lines) != 0 {
		t.Errorf("expected empty, got %v", lines)
	}
}

func TestFormatSkewWarnings_Content(t *testing.T) {
	ws := []SkewWarning{
		{Hour: 9, Minute: 0, Pattern: "09:00", Message: "2 jobs share clock time 09:00 — consider spreading minutes to reduce skew"},
	}
	lines := FormatSkewWarnings(ws)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] == "" {
		t.Error("expected non-empty line")
	}
}
