package analyzer

import (
	"fmt"
	"sort"
	"strings"
)

// Job represents a parsed cron job entry.
type DuplicateGroup struct {
	Expression string
	JobNames   []string
}

// DetectDuplicates finds jobs that share identical cron expressions.
// It returns a slice of DuplicateGroup, each containing the shared
// expression and the names of all jobs using it.
func DetectDuplicates(jobs []Job) []DuplicateGroup {
	exprIndex := make(map[string][]string)

	for _, job := range jobs {
		expr := normalizeExpression(job.Expression)
		exprIndex[expr] = append(exprIndex[expr], job.Name)
	}

	var groups []DuplicateGroup
	for expr, names := range exprIndex {
		if len(names) > 1 {
			sort.Strings(names)
			groups = append(groups, DuplicateGroup{
				Expression: expr,
				JobNames:   names,
			})
		}
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Expression < groups[j].Expression
	})

	return groups
}

// normalizeExpression collapses runs of whitespace so that
// "*  * * * *" and "* * * * *" are treated as identical.
func normalizeExpression(expr string) string {
	fields := strings.Fields(expr)
	return fmt.Sprintf("%s", strings.Join(fields, " "))
}
