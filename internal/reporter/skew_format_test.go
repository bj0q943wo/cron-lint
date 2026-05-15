package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/cron-lint/internal/analyzer"
)

func TestWriteSkewText_NoWarnings(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteSkewText(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No clock-skew") {
		t.Errorf("expected no-warning message, got: %q", buf.String())
	}
}

func TestWriteSkewText_WithWarnings(t *testing.T) {
	ws := []analyzer.SkewWarning{
		{
			Jobs:    []string{"alpha", "beta"},
			Hour:    9,
			Minute:  0,
			Pattern: "09:00",
			Message: "2 jobs share clock time 09:00 — consider spreading minutes to reduce skew",
		},
	}
	var buf bytes.Buffer
	if err := WriteSkewText(&buf, ws); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[SKEW]") {
		t.Errorf("expected [SKEW] badge, got: %q", out)
	}
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Errorf("expected job names in output, got: %q", out)
	}
	if !strings.Contains(out, "09:00") {
		t.Errorf("expected pattern 09:00 in output, got: %q", out)
	}
}

func TestWriteSkewJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteSkewJSON(&buf, []analyzer.SkewWarning{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &arr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(arr) != 0 {
		t.Errorf("expected empty array, got %d elements", len(arr))
	}
}

func TestWriteSkewJSON_WithWarnings(t *testing.T) {
	ws := []analyzer.SkewWarning{
		{
			Jobs:    []string{"job-a", "job-b"},
			Hour:    6,
			Minute:  15,
			Pattern: "06:15",
			Message: "2 jobs share clock time 06:15",
		},
	}
	var buf bytes.Buffer
	if err := WriteSkewJSON(&buf, ws); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &arr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(arr) != 1 {
		t.Fatalf("expected 1 element, got %d", len(arr))
	}
	if arr[0]["pattern"] != "06:15" {
		t.Errorf("unexpected pattern: %v", arr[0]["pattern"])
	}
	jobsRaw, ok := arr[0]["jobs"].([]interface{})
	if !ok || len(jobsRaw) != 2 {
		t.Errorf("expected 2 jobs in JSON, got %v", arr[0]["jobs"])
	}
}
