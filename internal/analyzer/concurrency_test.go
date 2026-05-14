package analyzer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func makeConcJob(name, expr string) Job {
	s, err := mustParseSchedule(expr)
	if err != nil {
		panic(err)
	}
	return Job{Name: name, Schedule: s}
}

func mustParseSchedule(expr string) (*Schedule, error) {
	return mustParseHelper(expr)
}

func TestCheckConcurrency_NoConflict(t *testing.T) {
	jobs := []Job{
		makeConcJob("job-a", "0 * * * *"),
		makeConcJob("job-b", "30 * * * *"),
	}
	got := CheckConcurrency(jobs)
	if len(got) != 0 {
		t.Fatalf("expected no warnings, got %d", len(got))
	}
}

func TestCheckConcurrency_ExactCollision(t *testing.T) {
	jobs := []Job{
		makeConcJob("alpha", "0 9 * * 1"),
		makeConcJob("beta", "0 9 * * 1"),
	}
	got := CheckConcurrency(jobs)
	if len(got) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(got))
	}
	if !containsAll(got[0].Jobs, []string{"alpha", "beta"}) {
		t.Errorf("unexpected jobs in warning: %v", got[0].Jobs)
	}
}

func TestCheckConcurrency_ThreeWayCollision(t *testing.T) {
	jobs := []Job{
		makeConcJob("x", "15 6 * * *"),
		makeConcJob("y", "15 6 * * *"),
		makeConcJob("z", "15 6 * * *"),
	}
	got := CheckConcurrency(jobs)
	// 7 days × 1 slot each = 7 warnings
	if len(got) == 0 {
		t.Fatal("expected warnings for three-way collision")
	}
	for _, w := range got {
		if len(w.Jobs) < 3 {
			t.Errorf("expected 3 jobs in collision, got %v", w.Jobs)
		}
	}
}

func TestCheckConcurrency_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{
		{Name: "nil-job", Schedule: nil},
		makeConcJob("real", "5 5 * * *"),
	}
	got := CheckConcurrency(jobs)
	if len(got) != 0 {
		t.Fatalf("expected no warnings, got %d", len(got))
	}
}

func TestCheckConcurrency_SuggestionNotEmpty(t *testing.T) {
	jobs := []Job{
		makeConcJob("p", "0 12 * * 3"),
		makeConcJob("q", "0 12 * * 3"),
	}
	got := CheckConcurrency(jobs)
	for _, w := range got {
		if w.Suggestion == "" {
			t.Error("expected non-empty suggestion")
		}
	}
	_ = cmp.Diff // keep import used
}

func containsAll(haystack, needles []string) bool {
	set := make(map[string]bool, len(haystack))
	for _, h := range haystack {
		set[h] = true
	}
	for _, n := range needles {
		if !set[n] {
			return false
		}
	}
	return true
}
