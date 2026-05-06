// Package reporter formats analysis results for human and machine consumers.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/cron-lint/internal/analyzer"
)

// Report holds the complete analysis output.
type Report struct {
	Jobs     []analyzer.Job              `json:"jobs"`
	Overlaps []analyzer.OverlapResult    `json:"overlaps"`
	Warnings []analyzer.ValidationWarning `json:"warnings,omitempty"`
}

// Build assembles a Report from the provided jobs, running overlap detection
// and schedule validation.
func Build(jobs []analyzer.Job) Report {
	return Report{
		Jobs:     jobs,
		Overlaps: analyzer.DetectOverlaps(jobs),
		Warnings: analyzer.ValidateJobs(jobs),
	}
}

// WriteText writes a human-readable report to w.
func WriteText(w io.Writer, r Report) error {
	fmt.Fprintf(w, "Jobs analysed: %d\n", len(r.Jobs))

	if len(r.Overlaps) == 0 {
		fmt.Fprintln(w, "No overlapping schedules detected.")
	} else {
		fmt.Fprintf(w, "Overlaps detected: %d\n", len(r.Overlaps))
		for _, o := range r.Overlaps {
			fmt.Fprintf(w, "  [OVERLAP] %q and %q share %d minute(s) — e.g. %s\n",
				o.JobA.Name, o.JobB.Name, len(o.CommonMinutes),
				formatSample(o.CommonMinutes))
		}
	}

	if len(r.Warnings) > 0 {
		fmt.Fprintf(w, "Warnings: %d\n", len(r.Warnings))
		for _, w2 := range r.Warnings {
			fmt.Fprintf(w, "  [WARN] %q: %s\n", w2.Job.Name, w2.Message)
		}
	}

	return nil
}

// WriteJSON writes a machine-readable JSON report to w.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// formatSample returns a short preview of minute values.
func formatSample(minutes []int) string {
	const max = 3
	if len(minutes) == 0 {
		return "(none)"
	}
	parts := make([]string, 0, max)
	for i, m := range minutes {
		if i >= max {
			parts = append(parts, "...")
			break
		}
		parts = append(parts, fmt.Sprintf("%d", m))
	}
	return strings.Join(parts, ",")
}
