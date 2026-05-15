package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/your-org/cron-lint/internal/analyzer"
)

// WriteSkewText writes skew warnings in human-readable form to w.
func WriteSkewText(w io.Writer, warnings []analyzer.SkewWarning) error {
	if len(warnings) == 0 {
		_, err := fmt.Fprintln(w, "No clock-skew warnings detected.")
		return err
	}
	for _, warn := range warnings {
		_, err := fmt.Fprintf(w, "[SKEW] %s  jobs: %s\n  → %s\n",
			warn.Pattern,
			strings.Join(warn.Jobs, ", "),
			warn.Message,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// skewWarningJSON is the wire representation for a single skew warning.
type skewWarningJSON struct {
	Pattern string   `json:"pattern"`
	Hour    int      `json:"hour"`
	Minute  int      `json:"minute"`
	Jobs    []string `json:"jobs"`
	Message string   `json:"message"`
}

// WriteSkewJSON writes skew warnings as a JSON array to w.
func WriteSkewJSON(w io.Writer, warnings []analyzer.SkewWarning) error {
	out := make([]skewWarningJSON, 0, len(warnings))
	for _, warn := range warnings {
		out = append(out, skewWarningJSON{
			Pattern: warn.Pattern,
			Hour:    warn.Hour,
			Minute:  warn.Minute,
			Jobs:    warn.Jobs,
			Message: warn.Message,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
