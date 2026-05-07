// Package reporter formats analysis results for human and machine consumption.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/cron-lint/internal/analyzer"
	"github.com/user/cron-lint/internal/parser"
)

// Report is the top-level result returned by Build.
type Report struct {
	Jobs        []parser.Job           `json:"jobs"`
	Duplicates  []analyzer.DuplicateGroup `json:"duplicates"`
	Overlaps    []analyzer.OverlapPair    `json:"overlaps"`
	Warnings    []analyzer.Warning        `json:"warnings"`
	Suggestions []analyzer.Suggestion     `json:"suggestions"`
}

// Build runs all analysers and assembles a Report.
func Build(jobs []parser.Job) Report {
	return Report{
		Jobs:        jobs,
		Duplicates:  analyzer.DetectDuplicates(jobs),
		Overlaps:    analyzer.DetectOverlaps(jobs),
		Warnings:    analyzer.ValidateJobs(jobs),
		Suggestions: analyzer.SuggestFixes(jobs),
	}
}

// WriteText writes a human-readable report to w.
func WriteText(w io.Writer, r Report) {
	fmt.Fprintf(w, "Jobs loaded: %d\n", len(r.Jobs))

	if len(r.Duplicates) == 0 && len(r.Overlaps) == 0 && len(r.Warnings) == 0 && len(r.Suggestions) == 0 {
		fmt.Fprintln(w, "No issues found.")
		return
	}

	if len(r.Duplicates) > 0 {
		fmt.Fprintf(w, "\nDuplicates (%d group(s)):\n", len(r.Duplicates))
		for _, d := range r.Duplicates {
			names := make([]string, len(d.Jobs))
			for i, j := range d.Jobs {
				names[i] = j.Name
			}
			fmt.Fprintf(w, "  [%s] share schedule %q\n", strings.Join(names, ", "), d.Schedule)
		}
	}

	if len(r.Overlaps) > 0 {
		fmt.Fprintf(w, "\nOverlaps (%d pair(s)):\n", len(r.Overlaps))
		for _, o := range r.Overlaps {
			fmt.Fprintf(w, "  %q and %q overlap at %s\n", o.A.Name, o.B.Name, formatSample(o.SampleMinutes))
		}
	}

	if len(r.Warnings) > 0 {
		fmt.Fprintf(w, "\nWarnings (%d):\n", len(r.Warnings))
		for _, w2 := range r.Warnings {
			fmt.Fprintf(w, "  [%s] %s\n", w2.JobName, w2.Message)
		}
	}

	if len(r.Suggestions) > 0 {
		fmt.Fprintf(w, "\nSuggestions (%d):\n", len(r.Suggestions))
		for _, s := range r.Suggestions {
			fmt.Fprintf(w, "  [%s] %s → try %q\n", s.JobName, s.Reason, s.Suggested)
		}
	}
}

// WriteJSON writes a JSON-encoded report to w.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func formatSample(minutes []int) string {
	if len(minutes) == 0 {
		return "(unknown)"
	}
	parts := make([]string, len(minutes))
	for i, m := range minutes {
		parts[i] = fmt.Sprintf("%d", m)
	}
	return "minutes " + strings.Join(parts, ",")
}
