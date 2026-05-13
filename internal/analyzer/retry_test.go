package analyzer

import (
	"strings"
	"testing"
)

func makeRetryJob(name, expr string) Job {
	sched, err := mustParseRetry(expr)
	if err != nil {
		panic(err)
	}
	return Job{Name: name, Expression: expr, Schedule: sched}
}

func mustParseRetry(expr string) (*Schedule, error) {
	return Parse(expr)
}

func TestCheckRetry_FrequentJob_NoWarning(t *testing.T) {
	// every minute — interval = 1
	job := makeRetryJob("frequent", "* * * * *")
	warnings := CheckRetry([]Job{job})
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

func TestCheckRetry_HourlyJob_NoWarning(t *testing.T) {
	// every hour — interval = 60
	job := makeRetryJob("hourly", "0 * * * *")
	warnings := CheckRetry([]Job{job})
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

func TestCheckRetry_WeeklyJob_Warning(t *testing.T) {
	// once a week — interval >> 1440
	job := makeRetryJob("weekly", "0 9 * * 1")
	warnings := CheckRetry([]Job{job})
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	w := warnings[0]
	if w.JobName != "weekly" {
		t.Errorf("unexpected job name %q", w.JobName)
	}
	if w.IntervalMinutes < retryThresholdMinutes {
		t.Errorf("interval %d should be >= %d", w.IntervalMinutes, retryThresholdMinutes)
	}
	if !strings.Contains(w.Message, "retry") {
		t.Errorf("message should mention retry: %s", w.Message)
	}
}

func TestCheckRetry_NilSchedule_Skipped(t *testing.T) {
	job := Job{Name: "broken", Expression: "bad", Schedule: nil}
	warnings := CheckRetry([]Job{job})
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings for nil schedule, got %d", len(warnings))
	}
}

func TestCheckRetry_MultipleJobs_OnlyLowFrequencyWarned(t *testing.T) {
	jobs := []Job{
		makeRetryJob("every-minute", "* * * * *"),
		makeRetryJob("daily", "0 3 * * *"),
		makeRetryJob("weekly", "0 9 * * 0"),
	}
	warnings := CheckRetry(jobs)
	// daily interval = 1440 minutes (exactly at threshold, not above) — no warning
	// weekly interval = 10080 — warning
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %+v", len(warnings), warnings)
	}
	if warnings[0].JobName != "weekly" {
		t.Errorf("expected warning for 'weekly', got %q", warnings[0].JobName)
	}
}

func TestFormatInterval(t *testing.T) {
	cases := []struct {
		minutes  int
		expected string
	}{
		{1, "1 minute(s)"},
		{90, "1 hour(s)"},
		{2880, "2 day(s)"},
		{20160, "2 week(s)"},
	}
	for _, c := range cases {
		got := formatInterval(c.minutes)
		if got != c.expected {
			t.Errorf("formatInterval(%d) = %q, want %q", c.minutes, got, c.expected)
		}
	}
}
