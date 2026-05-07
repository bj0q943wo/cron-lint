package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-org/cron-lint/internal/analyzer"
)

// WriteBurstText writes burst warnings in human-readable form to w.
func WriteBurstText(w io.Writer, warnings []analyzer.BurstWarning) {
	if len(warnings) == 0 {
		fmt.Fprintln(w, "[burst] no burst warnings detected")
		return
	}
	for _, warn := range warnings {
		fmt.Fprintf(w, "[burst] minute %02d — %d jobs fire within a %d-minute window\n",
			warn.Minute, warn.Count, warn.WindowMinutes)
		for _, job := range warn.Jobs {
			fmt.Fprintf(w, "         • %-20s  %s\n", job.Name, job.Expression)
		}
	}
}

type burstWarningJSON struct {
	StartMinute   int      `json:"start_minute"`
	WindowMinutes int      `json:"window_minutes"`
	Count         int      `json:"count"`
	Jobs          []string `json:"jobs"`
}

// WriteBurstJSON writes burst warnings as a JSON array to w.
func WriteBurstJSON(w io.Writer, warnings []analyzer.BurstWarning) error {
	out := make([]burstWarningJSON, 0, len(warnings))
	for _, warn := range warnings {
		names := make([]string, len(warn.Jobs))
		for i, j := range warn.Jobs {
			names[i] = j.Name
		}
		out = append(out, burstWarningJSON{
			StartMinute:   warn.Minute,
			WindowMinutes: warn.WindowMinutes,
			Count:         warn.Count,
			Jobs:          names,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
