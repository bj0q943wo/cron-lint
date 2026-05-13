package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/cron-lint/internal/analyzer"
)

func TestWriteDriftText_NoWarnings(t *testing.T) {
	var buf bytes.Buffer
	WriteDriftText(&buf, nil)
	if !strings.Contains(buf.String(), "No schedule drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteDriftText_WithWarnings(t *testing.T) {
	warnings := []analyzer.DriftWarning{
		{
			JobA:      "job-a",
			JobB:      "job-b",
			OffsetMin: 20,
			Message:   `jobs "job-a" and "job-b" share the same hour/day pattern but minute sets differ by 20 min (threshold: 5)`,
		},
	}
	var buf bytes.Buffer
	WriteDriftText(&buf, warnings)
	out := buf.String()
	if !strings.Contains(out, "1 drift warning") {
		t.Errorf("expected count line, got: %s", out)
	}
	if !strings.Contains(out, "job-a") || !strings.Contains(out, "job-b") {
		t.Errorf("expected job names in output, got: %s", out)
	}
	if !strings.Contains(out, "20 min") {
		t.Errorf("expected offset in output, got: %s", out)
	}
}

func TestWriteDriftJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteDriftJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty array, got %d elements", len(result))
	}
}

func TestWriteDriftJSON_WithWarnings(t *testing.T) {
	warnings := []analyzer.DriftWarning{
		{JobA: "alpha", JobB: "beta", OffsetMin: 10, Message: "drift detected"},
	}
	var buf bytes.Buffer
	if err := WriteDriftJSON(&buf, warnings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 element, got %d", len(result))
	}
	if result[0]["job_a"] != "alpha" {
		t.Errorf("expected job_a=alpha, got %v", result[0]["job_a"])
	}
	if result[0]["offset_minutes"].(float64) != 10 {
		t.Errorf("expected offset_minutes=10, got %v", result[0]["offset_minutes"])
	}
}
