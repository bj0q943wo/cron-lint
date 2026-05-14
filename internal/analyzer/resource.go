package analyzer

import (
	"fmt"
	"sort"

	"github.com/your-org/cron-lint/internal/parser"
)

// ResourceWarning describes a job whose schedule may saturate a shared resource.
type ResourceWarning struct {
	JobName  string
	Expr     string
	Slot     string // e.g. "00:05" — the minute slot that is crowded
	PeakJobs []string
	Message  string
}

// ResourceParams controls the thresholds used by CheckResourceContention.
type ResourceParams struct {
	// MaxJobsPerSlot is the maximum number of jobs allowed to fire in the same
	// minute slot before a warning is emitted. Defaults to 3.
	MaxJobsPerSlot int
}

// CheckResourceContention detects minute slots where too many jobs fire at
// once, which can saturate shared resources (DB connections, API rate-limits,
// etc.).
func CheckResourceContention(jobs []Job, params ResourceParams) []ResourceWarning {
	if params.MaxJobsPerSlot <= 0 {
		params.MaxJobsPerSlot = 3
	}

	// Build a map: "HH:MM" -> []jobName
	slotMap := make(map[string][]string)

	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}
		for _, h := range j.Schedule.Hours {
			for _, m := range j.Schedule.Minutes {
				key := fmt.Sprintf("%02d:%02d", h, m)
				slotMap[key] = append(slotMap[key], j.Name)
			}
		}
	}

	// Collect crowded slots
	type crowded struct {
		slot string
		names []string
	}
	var crowdedSlots []crowded
	for slot, names := range slotMap {
		if len(names) > params.MaxJobsPerSlot {
			crowdedSlots = append(crowdedSlots, crowded{slot, names})
		}
	}
	sort.Slice(crowdedSlots, func(i, j int) bool {
		return crowdedSlots[i].slot < crowdedSlots[j].slot
	})

	// Emit one warning per (job, slot) pair so callers can attribute warnings
	// back to individual jobs.
	var warnings []ResourceWarning
	for _, cs := range crowdedSlots {
		for _, name := range cs.names {
			expr := expressionForJob(jobs, name)
			warnings = append(warnings, ResourceWarning{
				JobName:  name,
				Expr:     expr,
				Slot:     cs.slot,
				PeakJobs: cs.names,
				Message: fmt.Sprintf(
					"slot %s has %d concurrent jobs (threshold %d): %v",
					cs.slot, len(cs.names), params.MaxJobsPerSlot, cs.names,
				),
			})
		}
	}
	return warnings
}

// expressionForJob returns the raw expression string for a named job.
func expressionForJob(jobs []Job, name string) string {
	for _, j := range jobs {
		if j.Name == name {
			return j.Raw
		}
	}
	return ""
}

// Ensure parser import is used transitively through the Job type.
var _ *parser.Schedule
