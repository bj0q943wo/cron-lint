package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/cron-lint/internal/analyzer"
)

// WriteConcurrencyText writes a human-readable summary of concurrency
// warnings to w. When there are no warnings a single "OK" line is emitted.
func WriteConcurrencyText(w io.Writer, warnings []analyzer.ConcurrencyWarning) {
	if len(warnings) == 0 {
		fmt.Fprintln(w, "concurrency: OK — no simultaneous job starts detected")
		return
	}
	// Deduplicate by sorted job-set so each unique collision group appears once.
	seen := make(map[string]bool)
	for _, warn := range warnings {
		copy := append([]string(nil), warn.Jobs...)
		sort.Strings(copy)
		key := strings.Join(copy, "|")
		if seen[key] {
			continue
		}
		seen[key] = true
		fmt.Fprintf(w, "[CONCURRENCY] jobs fire simultaneously: %s\n",
			strings.Join(copy, ", "))
		fmt.Fprintf(w, "  hint: %s\n", warn.Suggestion)
	}
}

// concurrencyJSONRecord is the JSON shape for a single concurrency warning.
type concurrencyJSONRecord struct {
	Jobs       []string `json:"jobs"`
	Suggestion string   `json:"suggestion"`
}

// WriteConcurrencyJSON writes warnings as a JSON array to w.
// An empty slice is serialised as `[]`.
func WriteConcurrencyJSON(w io.Writer, warnings []analyzer.ConcurrencyWarning) {
	// Collapse duplicate job-sets before serialising.
	seen := make(map[string]bool)
	var records []concurrencyJSONRecord
	for _, warn := range warnings {
		copy := append([]string(nil), warn.Jobs...)
		sort.Strings(copy)
		key := strings.Join(copy, "|")
		if seen[key] {
			continue
		}
		seen[key] = true
		records = append(records, concurrencyJSONRecord{
			Jobs:       copy,
			Suggestion: warn.Suggestion,
		})
	}
	if records == nil {
		records = []concurrencyJSONRecord{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(records)
}
