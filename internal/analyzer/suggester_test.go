package analyzer

import (
	"testing"

	"github.com/user/cron-lint/internal/parser"
)

func jobs(pairs ...string) []parser.Job {
	var out []parser.Job
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Job{Name: pairs[i], Schedule: pairs[i+1]})
	}
	return out
}

func TestSuggestFixes_NoSuggestions(t *testing.T) {
	input := jobs("backup", "30 2 * * *")
	got := SuggestFixes(input)
	if len(got) != 0 {
		t.Fatalf("expected 0 suggestions, got %d", len(got))
	}
}

func TestSuggestFixes_EveryMinute(t *testing.T) {
	input := jobs("poller", "* * * * *")
	got := SuggestFixes(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if got[0].Suggested != "*/5 * * * *" {
		t.Errorf("unexpected suggestion: %s", got[0].Suggested)
	}
}

func TestSuggestFixes_MidnightSpread(t *testing.T) {
	input := jobs("nightly", "0 0 * * *")
	got := SuggestFixes(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if got[0].Original != "0 0 * * *" {
		t.Errorf("original not preserved: %s", got[0].Original)
	}
	if got[0].Suggested == "0 0 * * *" {
		t.Error("suggestion should differ from original")
	}
}

func TestSuggestFixes_Multiple(t *testing.T) {
	input := jobs(
		"a", "* * * * *",
		"b", "0 0 * * 1",
		"c", "15 3 * * *",
	)
	got := SuggestFixes(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(got))
	}
}

func TestSuggestFixes_EmptyInput(t *testing.T) {
	got := SuggestFixes(nil)
	if got != nil && len(got) != 0 {
		t.Fatalf("expected no suggestions for empty input, got %d", len(got))
	}
}
