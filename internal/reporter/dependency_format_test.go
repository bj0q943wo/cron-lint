package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cron-lint/internal/analyzer"
)

func TestWriteDependencyText_NoWarnings(t *testing.T) {
	var buf bytes.Buffer
	WriteDependencyText(&buf, nil)
	if !strings.Contains(buf.String(), "no concerns") {
		t.Errorf("expected 'no concerns' message, got: %q", buf.String())
	}
}

func TestWriteDependencyText_WithWarnings(t *testing.T) {
	warnings := []analyzer.DependencyWarning{
		{JobA: "fetch", JobB: "process", Kind: "concurrent", Message: "jobs fire at same time"},
		{JobA: "extract", JobB: "transform", Kind: "successor", Message: "successor pattern detected"},
	}
	var buf bytes.Buffer
	WriteDependencyText(&buf, warnings)
	out := buf.String()
	if !strings.Contains(out, "2 concern(s)") {
		t.Errorf("expected count in output, got: %q", out)
	}
	if !strings.Contains(out, "CONCURRENT") {
		t.Errorf("expected CONCURRENT badge, got: %q", out)
	}
	if !strings.Contains(out, "SUCCESSOR") {
		t.Errorf("expected SUCCESSOR badge, got: %q", out)
	}
}

func TestWriteDependencyJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteDependencyJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var report DependencyReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if report.Total != 0 {
		t.Errorf("expected total=0, got %d", report.Total)
	}
	if len(report.Warnings) != 0 {
		t.Errorf("expected empty warnings slice, got %d", len(report.Warnings))
	}
}

func TestWriteDependencyJSON_WithWarnings(t *testing.T) {
	warnings := []analyzer.DependencyWarning{
		{JobA: "a", JobB: "b", Kind: "concurrent", Message: "overlap"},
	}
	var buf bytes.Buffer
	if err := WriteDependencyJSON(&buf, warnings); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var report DependencyReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if report.Total != 1 {
		t.Errorf("expected total=1, got %d", report.Total)
	}
	if report.Warnings[0].Kind != "concurrent" {
		t.Errorf("expected kind=concurrent, got %q", report.Warnings[0].Kind)
	}
}
