package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/cron-lint/internal/analyzer"
)

func makeJobs(exprs ...string) []analyzer.Job {
	var jobs []analyzer.Job
	for i, e := range exprs {
		jobs = append(jobs, analyzer.Job{
			Name:       fmt.Sprintf("job-%d", i+1),
			Expression: e,
		})
	}
	return jobs
}

func TestBuild_NoOverlaps(t *testing.T) {
	jobs := makeJobs("0 * * * *", "30 * * * *")
	r := Build(jobs)
	if len(r.Overlaps) != 0 {
		t.Errorf("expected no overlaps, got %d", len(r.Overlaps))
	}
	if len(r.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %d", len(r.Duplicates))
	}
}

func TestBuild_WithDuplicates(t *testing.T) {
	jobs := []analyzer.Job{
		{Name: "alpha", Expression: "*/5 * * * *"},
		{Name: "beta", Expression: "*/5 * * * *"},
		{Name: "gamma", Expression: "0 * * * *"},
	}
	r := Build(jobs)
	if len(r.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(r.Duplicates))
	}
	if r.Duplicates[0].Expression != "*/5 * * * *" {
		t.Errorf("unexpected duplicate expression %q", r.Duplicates[0].Expression)
	}
}

func TestWriteText_NoDuplicates(t *testing.T) {
	r := Report{}
	var buf bytes.Buffer
	WriteText(&buf, r)
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK message, got: %s", buf.String())
	}
}

func TestWriteText_WithDuplicates(t *testing.T) {
	r := Report{
		Duplicates: []analyzer.DuplicateGroup{
			{Expression: "0 * * * *", JobNames: []string{"job-1", "job-2"}},
		},
	}
	var buf bytes.Buffer
	WriteText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "DUPLICATES") {
		t.Errorf("expected DUPLICATES section, got: %s", out)
	}
	if !strings.Contains(out, "job-1") || !strings.Contains(out, "job-2") {
		t.Errorf("expected job names in output, got: %s", out)
	}
}

func TestWriteJSON_IncludesDuplicates(t *testing.T) {
	r := Report{
		Duplicates: []analyzer.DuplicateGroup{
			{Expression: "*/15 * * * *", JobNames: []string{"a", "b"}},
		},
	}
	var buf bytes.Buffer
	if err := WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var decoded Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(decoded.Duplicates) != 1 {
		t.Errorf("expected 1 duplicate in JSON output, got %d", len(decoded.Duplicates))
	}
}
