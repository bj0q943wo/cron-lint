package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cron-lint/internal/analyzer"
	"github.com/cron-lint/internal/reporter"
)

func makeJobs(n int) []analyzer.Job {
	jobs := make([]analyzer.Job, n)
	for i := range jobs {
		jobs[i] = analyzer.Job{Name: "job"}
	}
	return jobs
}

func TestBuild_NoOverlaps(t *testing.T) {
	r := reporter.Build(makeJobs(3), nil)
	if r.TotalJobs != 3 {
		t.Errorf("expected 3 jobs, got %d", r.TotalJobs)
	}
	if len(r.Overlaps) != 0 {
		t.Errorf("expected 0 overlaps, got %d", len(r.Overlaps))
	}
}

func TestBuild_WithOverlaps(t *testing.T) {
	overlaps := []analyzer.OverlapResult{
		{JobA: "backup", JobB: "report", SharedMinutes: []int{0, 30}},
	}
	r := reporter.Build(makeJobs(2), overlaps)
	if len(r.Overlaps) != 1 {
		t.Fatalf("expected 1 overlap report, got %d", len(r.Overlaps))
	}
	if r.Overlaps[0].JobA != "backup" || r.Overlaps[0].JobB != "report" {
		t.Errorf("unexpected job names in overlap report")
	}
	if !strings.Contains(r.Overlaps[0].Message, "2 overlapping minute") {
		t.Errorf("message missing minute count: %s", r.Overlaps[0].Message)
	}
}

func TestWriteText_NoOverlap(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.Report{TotalJobs: 2, Overlaps: nil}
	reporter.WriteText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "No overlapping") {
		t.Errorf("expected no-overlap message, got: %s", out)
	}
}

func TestWriteText_WithOverlap(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.Report{
		TotalJobs: 2,
		Overlaps: []reporter.OverlapReport{
			{JobA: "a", JobB: "b", Message: "jobs overlap"},
		},
	}
	reporter.WriteText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "Found 1 overlap") {
		t.Errorf("expected overlap count in output, got: %s", out)
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.Report{
		TotalJobs: 1,
		Overlaps: []reporter.OverlapReport{
			{JobA: "x", JobB: "y", Message: "overlap msg"},
		},
	}
	reporter.WriteJSON(&buf, r)
	out := buf.String()
	for _, want := range []string{"total_jobs", "overlap_count", "overlaps", "job_a", "job_b", "message"} {
		if !strings.Contains(out, want) {
			t.Errorf("JSON output missing key %q", want)
		}
	}
}
