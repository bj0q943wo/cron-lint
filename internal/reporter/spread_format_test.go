package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/cron-lint/internal/analyzer"
	"github.com/your-org/cron-lint/internal/parser"
)

func makeSpreadWarning(msg string, minutes []int, names ...string) analyzer.SpreadWarning {
	jobs := make([]analyzer.Job, len(names))
	for i, n := range names {
		sched, _ := parser.Parse("0 * * * *")
		jobs[i] = analyzer.Job{Name: n, Expression: "0 * * * *", Schedule: sched}
	}
	return analyzer.SpreadWarning{Jobs: jobs, Minutes: minutes, Message: msg}
}

func TestWriteSpreadText_NoWarnings(t *testing.T) {
	var buf bytes.Buffer
	WriteSpreadText(&buf, nil)
	if !strings.Contains(buf.String(), "No clustering") {
		t.Errorf("expected no-warning message, got: %s", buf.String())
	}
}

func TestWriteSpreadText_WithWarnings(t *testing.T) {
	w := makeSpreadWarning("4 jobs cluster within a 5-minute window starting at minute 0",
		[]int{0, 1, 2, 3}, "alpha", "beta", "gamma", "delta")
	var buf bytes.Buffer
	WriteSpreadText(&buf, []analyzer.SpreadWarning{w})
	out := buf.String()
	if !strings.Contains(out, "4 jobs cluster") {
		t.Errorf("expected cluster message in output, got: %s", out)
	}
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected job name 'alpha' in output")
	}
	if !strings.Contains(out, "[0 1 2 3]") {
		t.Errorf("expected minutes list in output, got: %s", out)
	}
}

func TestWriteSpreadJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteSpreadJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected empty array, got %d records", len(records))
	}
}

func TestWriteSpreadJSON_WithWarnings(t *testing.T) {
	w := makeSpreadWarning("3 jobs cluster", []int{5, 6, 7}, "x", "y", "z")
	var buf bytes.Buffer
	if err := WriteSpreadJSON(&buf, []analyzer.SpreadWarning{w}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	jobs, ok := records[0]["jobs"].([]interface{})
	if !ok || len(jobs) != 3 {
		t.Errorf("expected 3 jobs in JSON record")
	}
	mins, ok := records[0]["minutes"].([]interface{})
	if !ok || len(mins) != 3 {
		t.Errorf("expected 3 minutes in JSON record")
	}
}
