package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-org/cron-lint/internal/analyzer"
)

// WriteDriftText writes drift warnings in human-readable format to w.
func WriteDriftText(w io.Writer, warnings []analyzer.DriftWarning) {
	if len(warnings) == 0 {
		fmt.Fprintln(w, "[drift] No schedule drift detected.")
		return
	}
	fmt.Fprintf(w, "[drift] %d drift warning(s) found:\n", len(warnings))
	for _, warn := range warnings {
		fmt.Fprintf(w, "  ⚠  %s\n", warn.Message)
		fmt.Fprintf(w, "     jobs : %s  ↔  %s\n", warn.JobA, warn.JobB)
		fmt.Fprintf(w, "     offset: %d min\n", warn.OffsetMin)
	}
}

// driftJSON is the serialisable form of a DriftWarning.
type driftJSON struct {
	JobA      string `json:"job_a"`
	JobB      string `json:"job_b"`
	OffsetMin int    `json:"offset_minutes"`
	Message   string `json:"message"`
}

// WriteDriftJSON writes drift warnings as a JSON array to w.
func WriteDriftJSON(w io.Writer, warnings []analyzer.DriftWarning) error {
	out := make([]driftJSON, 0, len(warnings))
	for _, warn := range warnings {
		out = append(out, driftJSON{
			JobA:      warn.JobA,
			JobB:      warn.JobB,
			OffsetMin: warn.OffsetMin,
			Message:   warn.Message,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
