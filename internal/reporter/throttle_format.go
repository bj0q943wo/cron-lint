package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/your-org/cron-lint/internal/analyzer"
)

// WriteThrottleText writes human-readable throttle warnings to w.
func WriteThrottleText(w io.Writer, warnings []analyzer.ThrottleWarning) {
	if len(warnings) == 0 {
		fmt.Fprintln(w, "throttle: no high-frequency windows detected")
		return
	}
	// Sort by window for deterministic output.
	sorted := make([]analyzer.ThrottleWarning, len(warnings))
	copy(sorted, warnings)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Window < sorted[j].Window
	})

	fmt.Fprintf(w, "throttle: %d high-frequency window(s) detected\n", len(sorted))
	for _, warn := range sorted {
		sort.Strings(warn.Jobs)
		fmt.Fprintf(w, "  [THROTTLE] window %s — %d firings (threshold %d) jobs: %v\n",
			warn.Window, warn.Firings, warn.Threshold, warn.Jobs)
	}
}

// WriteThrottleJSON writes throttle warnings as a JSON array to w.
func WriteThrottleJSON(w io.Writer, warnings []analyzer.ThrottleWarning) error {
	type jsonWarning struct {
		Window    string   `json:"window"`
		Jobs      []string `json:"jobs"`
		Firings   int      `json:"firings"`
		Threshold int      `json:"threshold"`
	}

	out := make([]jsonWarning, 0, len(warnings))
	for _, w := range warnings {
		sorted := make([]string, len(w.Jobs))
		copy(sorted, w.Jobs)
		sort.Strings(sorted)
		out = append(out, jsonWarning{
			Window:    w.Window,
			Jobs:      sorted,
			Firings:   w.Firings,
			Threshold: w.Threshold,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Window < out[j].Window })

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
