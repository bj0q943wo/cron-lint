package analyzer

import (
	"testing"
)

func makeNamedJob(name, expr string) Job {
	sched := mustParse(expr)
	return Job{Name: name, Expression: expr, Schedule: sched}
}

func TestCheckDependencies_NoConcerns(t *testing.T) {
	jobs := []Job{
		makeNamedJob("backup", "0 2 * * *"),
		makeNamedJob("report", "0 6 * * *"),
	}
	warnings := CheckDependencies(jobs)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d: %+v", len(warnings), warnings)
	}
}

func TestCheckDependencies_ConcurrentJobs(t *testing.T) {
	jobs := []Job{
		makeNamedJob("fetch", "30 4 * * *"),
		makeNamedJob("process", "30 4 * * *"),
	}
	warnings := CheckDependencies(jobs)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Kind != "concurrent" {
		t.Errorf("expected kind=concurrent, got %q", warnings[0].Kind)
	}
	if warnings[0].JobA != "fetch" || warnings[0].JobB != "process" {
		t.Errorf("unexpected job names: %q %q", warnings[0].JobA, warnings[0].JobB)
	}
}

func TestCheckDependencies_SuccessorPattern(t *testing.T) {
	jobs := []Job{
		makeNamedJob("extract", "0 3 * * *"),
		makeNamedJob("transform", "1 3 * * *"),
	}
	warnings := CheckDependencies(jobs)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %+v", len(warnings), warnings)
	}
	if warnings[0].Kind != "successor" {
		t.Errorf("expected kind=successor, got %q", warnings[0].Kind)
	}
}

func TestCheckDependencies_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{
		{Name: "bad", Expression: "invalid", Schedule: nil},
		makeNamedJob("good", "0 1 * * *"),
	}
	warnings := CheckDependencies(jobs)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings for nil schedule, got %d", len(warnings))
	}
}

func TestCheckDependencies_MultipleOverlaps(t *testing.T) {
	jobs := []Job{
		makeNamedJob("a", "*/5 * * * *"),
		makeNamedJob("b", "*/5 * * * *"),
		makeNamedJob("c", "*/5 * * * *"),
	}
	warnings := CheckDependencies(jobs)
	// pairs: (a,b), (a,c), (b,c) => 3 concurrent warnings
	if len(warnings) != 3 {
		t.Fatalf("expected 3 warnings, got %d", len(warnings))
	}
	for _, w := range warnings {
		if w.Kind != "concurrent" {
			t.Errorf("expected concurrent, got %q", w.Kind)
		}
	}
}
