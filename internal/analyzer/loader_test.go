package analyzer_test

import (
	"strings"
	"testing"

	"github.com/cron-lint/internal/analyzer"
)

const validInput = `
# daily jobs
backup   0 2 * * *
report   30 6 * * 1-5

# frequent
heartbeat */5 * * * *
`

func TestLoadJobs_Valid(t *testing.T) {
	jobs, err := analyzer.LoadJobs(strings.NewReader(validInput))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 3 {
		t.Fatalf("expected 3 jobs, got %d", len(jobs))
	}
	if jobs[0].Name != "backup" {
		t.Errorf("expected first job to be 'backup', got %q", jobs[0].Name)
	}
	if jobs[1].Expression != "30 6 * * 1-5" {
		t.Errorf("unexpected expression: %q", jobs[1].Expression)
	}
}

func TestLoadJobs_InvalidFieldCount(t *testing.T) {
	input := "badjob * * *\n"
	_, err := analyzer.LoadJobs(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for wrong field count, got nil")
	}
}

func TestLoadJobs_InvalidExpression(t *testing.T) {
	input := "badjob 99 * * * *\n" // minute 99 is out of range
	_, err := analyzer.LoadJobs(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid cron expression, got nil")
	}
}

func TestLoadJobs_EmptyInput(t *testing.T) {
	jobs, err := analyzer.LoadJobs(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestLoadJobs_CommentsAndBlanks(t *testing.T) {
	input := "\n# comment\n   \n# another\n"
	jobs, err := analyzer.LoadJobs(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}
