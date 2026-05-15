package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/example/cron-lint/internal/analyzer"
)

// WriteIdleText writes idle-gap warnings in human-readable form to w.
func WriteIdleText(w io.Writer, warnings []analyzer.IdleWarning) error {
	if len(warnings) == 0 {
		_, err := fmt.Fprintln(w, "[idle] no idle gaps detected")
		return err
	}
	for _, warn := range warnings {
		_, err := fmt.Fprintf(w, "[idle] %s\n", warn.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

// idleWarningJSON is the JSON representation of an IdleWarning.
type idleWarningJSON struct {
	StartHour int    `json:"start_hour"`
	EndHour   int    `json:"end_hour"`
	GapHours  int    `json:"gap_hours"`
	Message   string `json:"message"`
}

// WriteIdleJSON writes idle-gap warnings as a JSON array to w.
func WriteIdleJSON(w io.Writer, warnings []analyzer.IdleWarning) error {
	out := make([]idleWarningJSON, 0, len(warnings))
	for _, warn := range warnings {
		out = append(out, idleWarningJSON{
			StartHour: warn.StartHour,
			EndHour:   warn.EndHour,
			GapHours:  warn.GapHours,
			Message:   warn.Message,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
