package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-org/cron-lint/internal/analyzer"
)

// WriteSpreadText writes spread warnings in human-readable form to w.
func WriteSpreadText(w io.Writer, warnings []analyzer.SpreadWarning) {
	if len(warnings) == 0 {
		fmt.Fprintln(w, "[spread] No clustering warnings.")
		return
	}
	fmt.Fprintf(w, "[spread] %d clustering warning(s) detected:\n", len(warnings))
	for _, warn := range warnings {
		fmt.Fprintf(w, "  ⚠  %s\n", warn.Message)
		fmt.Fprintf(w, "     minutes : %v\n", warn.Minutes)
		for _, j := range warn.Jobs {
			fmt.Fprintf(w, "       - %-20s  %s\n", j.Name, j.Expression)
		}
	}
}

// spreadJSONRecord is the JSON envelope for a single spread warning.
type spreadJSONRecord struct {
	Message string   `json:"message"`
	Minutes []int    `json:"minutes"`
	Jobs    []string `json:"jobs"`
}

// WriteSpreadJSON writes spread warnings as a JSON array to w.
func WriteSpreadJSON(w io.Writer, warnings []analyzer.SpreadWarning) error {
	records := make([]spreadJSONRecord, 0, len(warnings))
	for _, warn := range warnings {
		names := make([]string, len(warn.Jobs))
		for i, j := range warn.Jobs {
			names[i] = j.Name
		}
		records = append(records, spreadJSONRecord{
			Message: warn.Message,
			Minutes: warn.Minutes,
			Jobs:    names,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
