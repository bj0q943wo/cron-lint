package analyzer

import (
	"fmt"

	"github.com/your-org/cron-lint/internal/parser"
)

// SpreadWarning describes a cluster of jobs whose start minutes are too tightly
// packed within a single hour, increasing the risk of resource spikes.
type SpreadWarning struct {
	Jobs    []Job
	Minutes []int
	Message string
}

// SpreadOptions controls the behaviour of CheckSpread.
type SpreadOptions struct {
	// WindowSize is the rolling minute-window used to detect clusters (default 5).
	WindowSize int
	// MaxJobsInWindow is the maximum number of jobs allowed to start within
	// WindowSize minutes before a warning is emitted (default 3).
	MaxJobsInWindow int
}

// DefaultSpreadOptions returns sensible defaults for CheckSpread.
func DefaultSpreadOptions() SpreadOptions {
	return SpreadOptions{WindowSize: 5, MaxJobsInWindow: 3}
}

// CheckSpread detects jobs whose scheduled minutes cluster together within a
// rolling window, which can cause CPU / I/O spikes at those times.
func CheckSpread(jobs []Job, opts SpreadOptions) []SpreadWarning {
	if opts.WindowSize <= 0 {
		opts.WindowSize = 5
	}
	if opts.MaxJobsInWindow <= 0 {
		opts.MaxJobsInWindow = 3
	}

	// Build a map from minute-of-hour -> jobs that fire at that minute.
	byMinute := make(map[int][]Job)
	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}
		for _, m := range j.Schedule.Minutes {
			byMinute[m] = append(byMinute[m], j)
		}
	}

	var warnings []SpreadWarning
	seen := map[string]bool{}

	for start := 0; start < 60; start++ {
		var clusterJobs []Job
		var clusterMinutes []int
		for offset := 0; offset < opts.WindowSize; offset++ {
			m := (start + offset) % 60
			if js, ok := byMinute[m]; ok {
				clusterJobs = append(clusterJobs, js...)
				clusterMinutes = append(clusterMinutes, m)
			}
		}
		if len(clusterJobs) <= opts.MaxJobsInWindow {
			continue
		}
		// Deduplicate warnings by the set of job names.
		key := spreadKey(clusterJobs)
		if seen[key] {
			continue
		}
		seen[key] = true
		warnings = append(warnings, SpreadWarning{
			Jobs:    clusterJobs,
			Minutes: dedupInts(clusterMinutes),
			Message: fmt.Sprintf("%d jobs cluster within a %d-minute window starting at minute %d",
				len(clusterJobs), opts.WindowSize, start),
		})
	}
	return warnings
}

func spreadKey(jobs []Job) string {
	names := make([]string, len(jobs))
	for i, j := range jobs {
		names[i] = j.Name
	}
	return fmt.Sprintf("%v", names)
}

func dedupInts(in []int) []int {
	seen := map[int]bool{}
	var out []int
	for _, v := range in {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

// Ensure parser import is used transitively through Job.Schedule.
var _ = parser.Parse
