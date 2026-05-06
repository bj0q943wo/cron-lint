package analyzer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/cron-lint/internal/parser"
)

// LoadJobs reads a simple crontab-like format from r.
// Each non-blank, non-comment line must have the form:
//
//	<name> <min> <hour> <dom> <month> <dow>
func LoadJobs(r io.Reader) ([]Job, error) {
	var jobs []Job
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 6 {
			return nil, fmt.Errorf("line %d: expected <name> and 5 cron fields, got %d token(s)", lineNum, len(parts))
		}
		name := parts[0]
		expr := strings.Join(parts[1:], " ")
		sched, err := parser.Parse(expr)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid expression for job %q: %w", lineNum, name, err)
		}
		jobs = append(jobs, Job{
			Name:       name,
			Expression: expr,
			Schedule:   sched,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}
	return jobs, nil
}
