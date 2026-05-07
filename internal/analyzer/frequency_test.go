package analyzer

import (
	"testing"

	"github.com/user/cron-lint/internal/parser"
)

func makeJob(name, expr string) Job {
	sched, err := parser.Parse(expr)
	if err != nil {
		panic("makeJob: " + err.Error())
	}
	return Job{Name: name, Expression: expr, Schedule: sched}
}

func TestAnalyzeFrequency_EveryMinute(t *testing.T) {
	jobs := []Job{makeJob("noisy", "* * * * *")}
	reports := AnalyzeFrequency(jobs)
	if len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(reports))
	}
	r := reports[0]
	if r.RunsPerHour != 60 {
		t.Errorf("RunsPerHour: want 60, got %d", r.RunsPerHour)
	}
	if r.RunsPerDay != 60*24 {
		t.Errorf("RunsPerDay: want %d, got %d", 60*24, r.RunsPerDay)
	}
	if r.Category != "high" {
		t.Errorf("Category: want high, got %s", r.Category)
	}
}

func TestAnalyzeFrequency_EveryHour(t *testing.T) {
	jobs := []Job{makeJob("hourly", "0 * * * *")}
	reports := AnalyzeFrequency(jobs)
	r := reports[0]
	if r.RunsPerHour != 1 {
		t.Errorf("RunsPerHour: want 1, got %d", r.RunsPerHour)
	}
	if r.Category != "low" {
		t.Errorf("Category: want low, got %s", r.Category)
	}
}

func TestAnalyzeFrequency_Every10Minutes(t *testing.T) {
	jobs := []Job{makeJob("frequent", "*/10 * * * *")}
	reports := AnalyzeFrequency(jobs)
	r := reports[0]
	if r.RunsPerHour != 6 {
		t.Errorf("RunsPerHour: want 6, got %d", r.RunsPerHour)
	}
	if r.Category != "medium" {
		t.Errorf("Category: want medium, got %s", r.Category)
	}
}

func TestAnalyzeFrequency_MultipleJobs(t *testing.T) {
	jobs := []Job{
		makeJob("a", "0 0 * * *"),
		makeJob("b", "*/5 * * * *"),
	}
	reports := AnalyzeFrequency(jobs)
	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}
	if reports[0].Category != "low" {
		t.Errorf("job a: want low, got %s", reports[0].Category)
	}
	if reports[1].Category != "high" {
		t.Errorf("job b: want high, got %s", reports[1].Category)
	}
}

func TestFormatFrequency(t *testing.T) {
	r := FrequencyReport{
		Job:         makeJob("myjob", "0 * * * *"),
		RunsPerHour: 1,
		RunsPerDay:  24,
		Category:    "low",
	}
	out := FormatFrequency(r)
	if out == "" {
		t.Error("FormatFrequency returned empty string")
	}
}
