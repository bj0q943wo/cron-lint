// Package reporter formats and outputs analysis results for cron-lint.
package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/cron-lint/internal/analyzer"
)

// Format represents the output format for the reporter.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// OverlapReport holds a single overlap finding.
type OverlapReport struct {
	JobA    string
	JobB    string
	Message string
}

// Report holds the full analysis result.
type Report struct {
	TotalJobs int
	Overlaps  []OverlapReport
}

// Build constructs a Report from detected overlaps.
func Build(jobs []analyzer.Job, overlaps []analyzer.OverlapResult) Report {
	reports := make([]OverlapReport, 0, len(overlaps))
	for _, o := range overlaps {
		reports = append(reports, OverlapReport{
			JobA:    o.JobA,
			JobB:    o.JobB,
			Message: fmt.Sprintf("jobs %q and %q share %d overlapping minute(s) per hour", o.JobA, o.JobB, len(o.SharedMinutes)),
		})
	}
	return Report{
		TotalJobs: len(jobs),
		Overlaps:  reports,
	}
}

// WriteText writes a human-readable report to w.
func WriteText(w io.Writer, r Report) {
	fmt.Fprintf(w, "cron-lint: analyzed %d job(s)\n", r.TotalJobs)
	if len(r.Overlaps) == 0 {
		fmt.Fprintln(w, "No overlapping jobs detected.")
		return
	}
	fmt.Fprintf(w, "Found %d overlap(s):\n", len(r.Overlaps))
	for i, o := range r.Overlaps {
		fmt.Fprintf(w, "  [%d] %s\n", i+1, o.Message)
	}
}

// WriteJSON writes a JSON-formatted report to w.
func WriteJSON(w io.Writer, r Report) {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("  \"total_jobs\": %d,\n", r.TotalJobs))
	sb.WriteString(fmt.Sprintf("  \"overlap_count\": %d,\n", len(r.Overlaps)))
	sb.WriteString("  \"overlaps\": [\n")
	for i, o := range r.Overlaps {
		comma := ","
		if i == len(r.Overlaps)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    {\"job_a\": %q, \"job_b\": %q, \"message\": %q}%s\n",
			o.JobA, o.JobB, o.Message, comma))
	}
	sb.WriteString("  ]\n}\n")
	fmt.Fprint(w, sb.String())
}
