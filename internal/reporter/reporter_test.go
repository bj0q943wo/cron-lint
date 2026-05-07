package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/cron-lint/internal/parser"
)

func makeJobs(pairs ...string) []parser.Job {
	var out []parser.Job
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Job{Name: pairs[i], Schedule: pairs[i+1]})
	}
	return out
}

func TestBuild_NoOverlaps(t *testing.T) {
	jobs := makeJobs("a", "0 1 * * *", "b", "0 2 * * *")
	r := Build(jobs)
	if len(r.Overlaps) != 0 {
		t.Errorf("expected 0 overlaps, got %d", len(r.Overlaps))
	}
	if len(r.Jobs) != 2 {
		t.Errorf("expected 2 jobs in report")
	}
}

func TestBuild_WithDuplicates(t *testing.T) {
	jobs := makeJobs("a", "0 1 * * *", "b", "0 1 * * *")
	r := Build(jobs)
	if len(r.Duplicates) != 1 {
		t.Errorf("expected 1 duplicate group, got %d", len(r.Duplicates))
	}
}

func TestBuild_WithSuggestions(t *testing.T) {
	jobs := makeJobs("poll", "* * * * *")
	r := Build(jobs)
	if len(r.Suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(r.Suggestions))
	}
}

func TestWriteText_NoDuplicates(t *testing.T) {
	jobs := makeJobs("a", "0 1 * * *")
	r := Build(jobs)
	var buf bytes.Buffer
	WriteText(&buf, r)
	if !strings.Contains(buf.String(), "No issues found") {
		t.Errorf("expected 'No issues found', got: %s", buf.String())
	}
}

func TestWriteText_WithSuggestions(t *testing.T) {
	jobs := makeJobs("poll", "* * * * *")
	r := Build(jobs)
	var buf bytes.Buffer
	WriteText(&buf, r)
	if !strings.Contains(buf.String(), "Suggestions") {
		t.Errorf("expected Suggestions section, got: %s", buf.String())
	}
}

func TestWriteJSON_Valid(t *testing.T) {
	jobs := makeJobs("a", "0 1 * * *")
	r := Build(jobs)
	var buf bytes.Buffer
	if err := WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["suggestions"]; !ok {
		t.Error("JSON output missing 'suggestions' key")
	}
}
