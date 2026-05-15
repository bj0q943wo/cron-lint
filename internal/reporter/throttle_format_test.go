package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/cron-lint/internal/analyzer"
)

func TestWriteThrottleText_NoWarnings(t *testing.T) {
	var buf bytes.Buffer
	WriteThrottleText(&buf, nil)
	if !strings.Contains(buf.String(), "no high-frequency") {
		t.Errorf("expected no-warning message, got: %q", buf.String())
	}
}

func TestWriteThrottleText_WithWarnings(t *testing.T) {
	warnings := []analyzer.ThrottleWarning{
		{Window: "01:00-01:04", Jobs: []string{"b", "a"}, Firings: 8, Threshold: 5},
		{Window: "00:00-00:04", Jobs: []string{"c"}, Firings: 6, Threshold: 5},
	}
	var buf bytes.Buffer
	WriteThrottleText(&buf, warnings)
	out := buf.String()
	if !strings.Contains(out, "2 high-frequency") {
		t.Errorf("expected count in header, got: %q", out)
	}
	// Sorted by window: 00: before 01:
	if idx00 := strings.Index(out, "00:00"); idx00 == -1 {
		t.Error("expected 00:00 window in output")
	}
	if idx01 := strings.Index(out, "01:00"); idx01 == -1 {
		t.Error("expected 01:00 window in output")
	}
	if !strings.Contains(out, "[THROTTLE]") {
		t.Error("expected [THROTTLE] badge")
	}
}

func TestWriteThrottleJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteThrottleJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var arr []interface{}
	if err := json.Unmarshal(buf.Bytes(), &arr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(arr) != 0 {
		t.Errorf("expected empty array, got %d elements", len(arr))
	}
}

func TestWriteThrottleJSON_WithWarnings(t *testing.T) {
	warnings := []analyzer.ThrottleWarning{
		{Window: "03:00-03:04", Jobs: []string{"x", "y"}, Firings: 12, Threshold: 10},
	}
	var buf bytes.Buffer
	if err := WriteThrottleJSON(&buf, warnings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &arr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(arr) != 1 {
		t.Fatalf("expected 1 element, got %d", len(arr))
	}
	if arr[0]["window"] != "03:00-03:04" {
		t.Errorf("unexpected window: %v", arr[0]["window"])
	}
	if int(arr[0]["firings"].(float64)) != 12 {
		t.Errorf("unexpected firings: %v", arr[0]["firings"])
	}
}
