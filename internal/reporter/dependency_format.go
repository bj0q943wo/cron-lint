package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cron-lint/internal/analyzer"
)

// DependencyReport holds the formatted output for dependency warnings.
type DependencyReport struct {
	Total    int                           `json:"total"`
	Warnings []analyzer.DependencyWarning `json:"warnings"`
}

// WriteDependencyText writes human-readable dependency warnings to w.
func WriteDependencyText(w io.Writer, warnings []analyzer.DependencyWarning) {
	if len(warnings) == 0 {
		fmt.Fprintln(w, "dependency: no concerns detected")
		return
	}
	fmt.Fprintf(w, "dependency: %d concern(s) detected\n", len(warnings))
	for _, warn := range warnings {
		kindLabel := kindBadge(warn.Kind)
		fmt.Fprintf(w, "  [%s] %s\n", kindLabel, warn.Message)
	}
}

// WriteDependencyJSON writes dependency warnings as a JSON object to w.
func WriteDependencyJSON(w io.Writer, warnings []analyzer.DependencyWarning) error {
	report := DependencyReport{
		Total:    len(warnings),
		Warnings: warnings,
	}
	if report.Warnings == nil {
		report.Warnings = []analyzer.DependencyWarning{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func kindBadge(kind string) string {
	switch kind {
	case "concurrent":
		return "CONCURRENT"
	case "successor":
		return "SUCCESSOR "
	default:
		return "UNKNOWN   "
	}
}
