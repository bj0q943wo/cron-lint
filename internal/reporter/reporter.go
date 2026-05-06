// Package reporter formats analysis results for human and machine consumption.
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
	Overlaps   []analyzer.OverlapResult   `json:"overlaps"`
	Warnings   []analyzer.Warning         `json:"warnings"`
	Duplicates []analyzer.DuplicateGroup  `json:"duplicates"`
}

// Build constructs a Report from the provided jobs.
func Build(jobs []analyzer.Job) Report {
	return Report{
		Overlaps:   analyzer.DetectOverlaps(jobs),
		Warnings:   analyzer.ValidateJobs(jobs),
		Duplicates: analyzer.DetectDuplicates(jobs),
	}
}

// WriteText writes a human-readable report to w.
func WriteText(w io.Writer, r Report) {
	if len(r.Overlaps) == 0 && len(r.Warnings) == 0 && len(r.Duplicates) == 0 {
		fmt.Fprintln(w, "OK: no issues found.")
		return
	}

	if len(r.Overlaps) > 0 {
		fmt.Fprintf(w, "OVERLAPS (%d):\n", len(r.Overlaps))
		for _, o := range r.Overlaps {
			fmt.Fprintf(w, "  [overlap] %q and %q share %d firing time(s), e.g. %s\n",
				o.JobA, o.JobB, len(o.CommonMinutes), formatSample(o.CommonMinutes))
		}
	}

	if len(r.Warnings) > 0 {
		fmt.Fprintf(w, "WARNINGS (%d):\n", len(r.Warnings))
		for _, w2 := range r.Warnings {
			fmt.Fprintf(w, "  [warn] %s: %s\n", w2.JobName, w2.Message)
		}
	}

	if len(r.Duplicates) > 0 {
		fmt.Fprintf(w, "DUPLICATES (%d):\n", len(r.Duplicates))
		for _, d := range r.Duplicates {
			fmt.Fprintf(w, "  [duplicate] expression %q used by: %s\n",
				d.Expression, strings.Join(d.JobNames, ", "))
		}
	}
}

// WriteJSON writes a JSON-encoded report to w.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// formatSample returns a short string representation of the first minute value.
func formatSample(minutes []int) string {
	if len(minutes) == 0 {
		return "(none)"
	}
	return fmt.Sprintf("minute %d", minutes[0])
}
