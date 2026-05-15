package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/cron-lint/internal/analyzer"
)

func TestWriteIdleText_NoWarnings(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteIdleText(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no idle gaps") {
		t.Errorf("expected 'no idle gaps' message, got: %q", buf.String())
	}
}

func TestWriteIdleText_WithWarnings(t *testing.T) {
	warnings := []analyzer.IdleWarning{
		{StartHour: 2, EndHour: 7, GapHours: 5, Message: "no jobs scheduled for 5 consecutive hours (02:00–07:00)"},
	}
	var buf bytes.Buffer
	if err := WriteIdleText(&buf, warnings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "[idle]") {
		t.Errorf("expected '[idle]' prefix, got: %q", got)
	}
	if !strings.Contains(got, "02:00") {
		t.Errorf("expected hour in output, got: %q", got)
	}
}

func TestWriteIdleJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteIdleJSON(&buf, []analyzer.IdleWarning{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty array, got %d elements", len(out))
	}
}

func TestWriteIdleJSON_WithWarnings(t *testing.T) {
	warnings := []analyzer.IdleWarning{
		{StartHour: 10, EndHour: 14, GapHours: 4, Message: "no jobs scheduled for 4 consecutive hours (10:00–14:00)"},
		{StartHour: 20, EndHour: 23, GapHours: 3, Message: "no jobs scheduled for 3 consecutive hours (20:00–23:00)"},
	}
	var buf bytes.Buffer
	if err := WriteIdleJSON(&buf, warnings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if int(out[0]["gap_hours"].(float64)) != 4 {
		t.Errorf("expected gap_hours=4, got %v", out[0]["gap_hours"])
	}
	if out[1]["message"] == "" {
		t.Error("expected non-empty message in second entry")
	}
}
