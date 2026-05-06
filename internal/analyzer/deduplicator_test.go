package analyzer

import (
	"testing"
)

func TestDetectDuplicates_NoDuplicates(t *testing.T) {
	jobs := []Job{
		{Name: "job-a", Expression: "0 * * * *"},
		{Name: "job-b", Expression: "*/5 * * * *"},
		{Name: "job-c", Expression: "0 9 * * 1"},
	}

	groups := DetectDuplicates(jobs)
	if len(groups) != 0 {
		t.Fatalf("expected no duplicate groups, got %d", len(groups))
	}
}

func TestDetectDuplicates_ExactDuplicate(t *testing.T) {
	jobs := []Job{
		{Name: "job-a", Expression: "0 * * * *"},
		{Name: "job-b", Expression: "0 * * * *"},
		{Name: "job-c", Expression: "*/5 * * * *"},
	}

	groups := DetectDuplicates(jobs)
	if len(groups) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(groups))
	}
	if groups[0].Expression != "0 * * * *" {
		t.Errorf("unexpected expression %q", groups[0].Expression)
	}
	if len(groups[0].JobNames) != 2 {
		t.Errorf("expected 2 job names, got %d", len(groups[0].JobNames))
	}
}

func TestDetectDuplicates_WhitespaceNormalization(t *testing.T) {
	jobs := []Job{
		{Name: "job-a", Expression: "0  *  *  *  *"},
		{Name: "job-b", Expression: "0 * * * *"},
	}

	groups := DetectDuplicates(jobs)
	if len(groups) != 1 {
		t.Fatalf("expected 1 duplicate group after normalization, got %d", len(groups))
	}
}

func TestDetectDuplicates_MultipleGroups(t *testing.T) {
	jobs := []Job{
		{Name: "alpha", Expression: "*/10 * * * *"},
		{Name: "beta", Expression: "0 0 * * *"},
		{Name: "gamma", Expression: "*/10 * * * *"},
		{Name: "delta", Expression: "0 0 * * *"},
		{Name: "epsilon", Expression: "0 0 * * *"},
	}

	groups := DetectDuplicates(jobs)
	if len(groups) != 2 {
		t.Fatalf("expected 2 duplicate groups, got %d", len(groups))
	}

	// groups are sorted by expression
	if groups[0].Expression != "*/10 * * * *" {
		t.Errorf("unexpected first group expression %q", groups[0].Expression)
	}
	if len(groups[1].JobNames) != 3 {
		t.Errorf("expected 3 jobs in second group, got %d", len(groups[1].JobNames))
	}
}

func TestDetectDuplicates_EmptyInput(t *testing.T) {
	groups := DetectDuplicates(nil)
	if len(groups) != 0 {
		t.Fatalf("expected no groups for empty input, got %d", len(groups))
	}
}
