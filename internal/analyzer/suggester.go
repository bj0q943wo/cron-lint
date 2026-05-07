package analyzer

import (
	"fmt"
	"strings"

	"github.com/user/cron-lint/internal/parser"
)

// Suggestion holds a recommended alternative cron expression.
type Suggestion struct {
	JobName    string
	Original   string
	Suggested  string
	Reason     string
}

// SuggestFixes analyses jobs and returns human-readable suggestions for
// common scheduling anti-patterns (e.g. "every minute" jobs, midnight pile-ups).
func SuggestFixes(jobs []parser.Job) []Suggestion {
	var suggestions []Suggestion

	for _, job := range jobs {
		if s, ok := suggestEveryMinute(job); ok {
			suggestions = append(suggestions, s)
			continue
		}
		if s, ok := suggestMidnightSpread(job); ok {
			suggestions = append(suggestions, s)
		}
	}
	return suggestions
}

// suggestEveryMinute detects "* * * * *" and recommends a less aggressive schedule.
func suggestEveryMinute(job parser.Job) (Suggestion, bool) {
	parts := strings.Fields(job.Schedule)
	if len(parts) != 5 {
		return Suggestion{}, false
	}
	for _, p := range parts {
		if p != "*" {
			return Suggestion{}, false
		}
	}
	return Suggestion{
		JobName:   job.Name,
		Original:  job.Schedule,
		Suggested: "*/5 * * * *",
		Reason:    "running every minute is rarely necessary; consider every 5 minutes",
	}, true
}

// suggestMidnightSpread detects jobs pinned to minute 0 and hour 0 and
// recommends spreading them to avoid thundering-herd at midnight.
func suggestMidnightSpread(job parser.Job) (Suggestion, bool) {
	parts := strings.Fields(job.Schedule)
	if len(parts) != 5 {
		return Suggestion{}, false
	}
	minute, hour := parts[0], parts[1]
	if minute != "0" || hour != "0" {
		return Suggestion{}, false
	}
	// Suggest a pseudo-random offset derived from the job name length.
	offset := len(job.Name) % 59
	if offset == 0 {
		offset = 7
	}
	suggested := fmt.Sprintf("%d 0 %s %s %s", offset, parts[2], parts[3], parts[4])
	return Suggestion{
		JobName:   job.Name,
		Original:  job.Schedule,
		Suggested: suggested,
		Reason:    fmt.Sprintf("many jobs fire at midnight; consider offsetting by %d minute(s)", offset),
	}, true
}
