package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// CronExpression represents a parsed cron schedule.
type CronExpression struct {
	Raw     string
	Minute  Field
	Hour    Field
	Day     Field
	Month   Field
	Weekday Field
}

// Field represents a single cron field with its resolved values.
type Field struct {
	Raw    string
	Values []int
}

var fieldRanges = [5][2]int{
	{0, 59}, // minute
	{0, 23}, // hour
	{1, 31}, // day
	{1, 12}, // month
	{0, 6},  // weekday
}

// Parse parses a standard 5-field cron expression.
func Parse(expr string) (*CronExpression, error) {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	fields := make([]Field, 5)
	for i, part := range parts {
		values, err := expandField(part, fieldRanges[i][0], fieldRanges[i][1])
		if err != nil {
			return nil, fmt.Errorf("field %d (%q): %w", i+1, part, err)
		}
		fields[i] = Field{Raw: part, Values: values}
	}

	return &CronExpression{
		Raw:     expr,
		Minute:  fields[0],
		Hour:    fields[1],
		Day:     fields[2],
		Month:   fields[3],
		Weekday: fields[4],
	}, nil
}

func expandField(field string, min, max int) ([]int, error) {
	if field == "*" {
		return rangeSlice(min, max, 1), nil
	}

	var result []int
	for _, part := range strings.Split(field, ",") {
		vals, err := expandPart(part, min, max)
		if err != nil {
			return nil, err
		}
		result = append(result, vals...)
	}
	return deduplicate(result), nil
}

func expandPart(part string, min, max int) ([]int, error) {
	if strings.Contains(part, "/") {
		sub := strings.SplitN(part, "/", 2)
		step, err := strconv.Atoi(sub[1])
		if err != nil || step <= 0 {
			return nil, fmt.Errorf("invalid step %q", sub[1])
		}
		base, err := expandPart(sub[0], min, max)
		if err != nil {
			return nil, err
		}
		var out []int
		for i, v := range base {
			if i%step == 0 {
				out = append(out, v)
			}
		}
		return out, nil
	}
	if strings.Contains(part, "-") {
		sub := strings.SplitN(part, "-", 2)
		lo, err1 := strconv.Atoi(sub[0])
		hi, err2 := strconv.Atoi(sub[1])
		if err1 != nil || err2 != nil || lo > hi || lo < min || hi > max {
			return nil, fmt.Errorf("invalid range %q", part)
		}
		return rangeSlice(lo, hi, 1), nil
	}
	if part == "*" {
		return rangeSlice(min, max, 1), nil
	}
	n, err := strconv.Atoi(part)
	if err != nil || n < min || n > max {
		return nil, fmt.Errorf("value %q out of range [%d,%d]", part, min, max)
	}
	return []int{n}, nil
}

func rangeSlice(lo, hi, step int) []int {
	var out []int
	for i := lo; i <= hi; i += step {
		out = append(out, i)
	}
	return out
}

func deduplicate(vals []int) []int {
	seen := make(map[int]struct{})
	var out []int
	for _, v := range vals {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
