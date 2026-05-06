package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/cron-lint/internal/analyzer"
	"github.com/user/cron-lint/internal/parser"
)

func makeJobs(exprs map[string]string) []analyzer.Job {
	var jobs []analyzer.Job
	for name, expr := range exprs {
		sched, err := parser.Parse(expr)
		if err != nil {
			panic(err)
		}
		jobs = append(jobs, analyzer.Job{Name: name, Expr: expr, Schedule: sched})
	}
	return jobs
}

func TestBuild_NoOverlaps(t *testing.T) {
	jobs := makeJobs(map[string]string{
		"a": "0 1 * * *",
		"b": "0 2 * * *",
	})
	r := Build(jobs)
	if len(r.Overlaps) != 0 {
		t.Errorf("expected 0 overlaps, got %d", len(r.Overlaps))
	}
}

func TestBuild_WithOverlaps(t *testing.T) {
	jobs := makeJobs(map[string]string{
		"x": "0 * * * *",
		"y": "0 * * * *",
	})
	r := Build(jobs)
	if len(r.Overlaps) == 0 {
		t.Error("expected overlaps, got none")
	}
}

func TestBuild_WithWarnings(t *testing.T) {
	jobs := makeJobs(map[string]string{
		"poller": "* * * * *",
	})
	r := Build(jobs)
	if len(r.Warnings) == 0 {
		t.Error("expected at least one warning for every-minute job")
	}
}

func TestWriteText_NoOverlap(t *testing.T) {
	jobs := makeJobs(map[string]string{"j": "5 4 * * *"})
	r := Build(jobs)
	var buf bytes.Buffer
	if err := WriteText(&buf, r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No overlapping") {
		t.Errorf("expected 'No overlapping' in output, got: %s", buf.String())
	}
}

func TestWriteText_WithOverlap(t *testing.T) {
	jobs := makeJobs(map[string]string{
		"a": "0 6 * * *",
		"b": "0 6 * * *",
	})
	r := Build(jobs)
	var buf bytes.Buffer
	if err := WriteText(&buf, r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "OVERLAP") {
		t.Errorf("expected OVERLAP in output, got: %s", buf.String())
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	jobs := makeJobs(map[string]string{"nightly": "0 0 * * *"})
	r := Build(jobs)
	var buf bytes.Buffer
	if err := WriteJSON(&buf, r); err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["jobs"]; !ok {
		t.Error("JSON output missing 'jobs' key")
	}
	if _, ok := out["overlaps"]; !ok {
		t.Error("JSON output missing 'overlaps' key")
	}
}
