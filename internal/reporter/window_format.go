package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-org/cron-lint/internal/analyzer"
)

// WriteWindowText writes a human-readable report of window warnings to w.
func WriteWindowText(w io.Writer, warnings []analyzer.WindowWarning) {
	if len(warnings) == 0 {
		fmt.Fprintln(w, "No window conflicts detected.")
		return
	}
	fmt.Fprintf(w, "Window Conflicts (%d):\n", len(warnings))
	for _, warn := range warnings {
		fmt.Fprintf(w, "  [WINDOW] %s\n", warn.Message)
		fmt.Fprintf(w, "           job: %s  expr: %s\n",
			warn.Job.Name, warn.Job.Expression)
	}
}

type jsonWindowWarning struct {
	Job     string `json:"job"`
	Expr    string `json:"expression"`
	Window  string `json:"window"`
	Start   int    `json:"window_start_min"`
	End     int    `json:"window_end_min"`
	Message string `json:"message"`
}

// WriteWindowJSON writes window warnings as a JSON array to w.
func WriteWindowJSON(w io.Writer, warnings []analyzer.WindowWarning) error {
	out := make([]jsonWindowWarning, 0, len(warnings))
	for _, warn := range warnings {
		out = append(out, jsonWindowWarning{
			Job:     warn.Job.Name,
			Expr:    warn.Job.Expression,
			Window:  warn.Window.Name,
			Start:   warn.Window.Start,
			End:     warn.Window.End,
			Message: warn.Message,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
