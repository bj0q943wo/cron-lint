package parser

import (
	"testing"
)

func TestParse_ValidExpressions(t *testing.T) {
	tests := []struct {
		expr        string
		minuteLen   int
		hourLen     int
	}{
		{"* * * * *", 60, 24},
		{"0 * * * *", 1, 24},
		{"0 12 * * *", 1, 1},
		{"*/15 * * * *", 4, 24},
		{"0-5 * * * *", 6, 24},
		{"0,30 9-17 * * 1-5", 2, 9},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			ce, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(ce.Minute.Values) != tt.minuteLen {
				t.Errorf("minute values: got %d, want %d", len(ce.Minute.Values), tt.minuteLen)
			}
			if len(ce.Hour.Values) != tt.hourLen {
				t.Errorf("hour values: got %d, want %d", len(ce.Hour.Values), tt.hourLen)
			}
			if ce.Raw != tt.expr {
				t.Errorf("Raw: got %q, want %q", ce.Raw, tt.expr)
			}
		})
	}
}

func TestParse_InvalidExpressions(t *testing.T) {
	tests := []struct {
		expr string
	}{
		{"* * * *"},           // too few fields
		{"* * * * * *"},       // too many fields
		{"60 * * * *"},        // minute out of range
		{"* 24 * * *"},        // hour out of range
		{"* * 0 * *"},         // day out of range
		{"* * * 13 *"},        // month out of range
		{"* * * * 7"},         // weekday out of range
		{"abc * * * *"},       // non-numeric
		{"5-3 * * * *"},       // invalid range
		{"*/0 * * * *"},       // zero step
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			_, err := Parse(tt.expr)
			if err == nil {
				t.Errorf("expected error for %q, got nil", tt.expr)
			}
		})
	}
}

func TestExpandField_StepOnWildcard(t *testing.T) {
	vals, err := expandField("*/10", 0, 59)
	if err != nil {
		t.Fatal(err)
	}
	expected := []int{0, 10, 20, 30, 40, 50}
	if len(vals) != len(expected) {
		t.Fatalf("got %v, want %v", vals, expected)
	}
	for i, v := range vals {
		if v != expected[i] {
			t.Errorf("vals[%d] = %d, want %d", i, v, expected[i])
		}
	}
}

func TestExpandField_CommaList(t *testing.T) {
	vals, err := expandField("1,15,30,45", 0, 59)
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 4 {
		t.Errorf("expected 4 values, got %d: %v", len(vals), vals)
	}
}
