package analyzer

import (
	"testing"
)

func makeJitterJob(name, raw string, minutes []int) Job {
	return Job{
		Name: name,
		Raw:  raw,
		Schedule: &Schedule{
			Minutes: minutes,
			Hours:   []int{0, 6, 12, 18},
			Days:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
			Months:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Weekdays: []int{0, 1, 2, 3, 4, 5, 6},
		},
	}
}

func TestCheckJitter_NoConflict(t *testing.T) {
	jobs := []Job{
		makeJitterJob("a", "0 * * * *", []int{0}),
		makeJitterJob("b", "15 * * * *", []int{15}),
		makeJitterJob("c", "30 * * * *", []int{30}),
	}
	warnings := CheckJitter(jobs)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

func TestCheckJitter_TwoJobsSameMinute(t *testing.T) {
	jobs := []Job{
		makeJitterJob("alpha", "0 * * * *", []int{0}),
		makeJitterJob("beta", "0 * * * *", []int{0}),
	}
	warnings := CheckJitter(jobs)
	if len(warnings) != 2 {
		t.Fatalf("expected 2 warnings (one per job), got %d", len(warnings))
	}
	for _, w := range warnings {
		if w.Suggest == "" {
			t.Errorf("expected a suggestion, got empty string for job %s", w.Job.Name)
		}
	}
}

func TestCheckJitter_ThreeJobsSameMinute(t *testing.T) {
	jobs := []Job{
		makeJitterJob("x", "5 * * * *", []int{5}),
		makeJitterJob("y", "5 * * * *", []int{5}),
		makeJitterJob("z", "5 * * * *", []int{5}),
	}
	warnings := CheckJitter(jobs)
	if len(warnings) != 3 {
		t.Fatalf("expected 3 warnings, got %d", len(warnings))
	}
}

func TestCheckJitter_NilScheduleSkipped(t *testing.T) {
	jobs := []Job{
		{Name: "no-sched", Raw: "* * * * *", Schedule: nil},
		makeJitterJob("real", "0 * * * *", []int{0}),
	}
	warnings := CheckJitter(jobs)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings when only one valid schedule, got %d", len(warnings))
	}
}

func TestCheckJitter_SuggestOffsetWrapsAround(t *testing.T) {
	// minute 55 + 7 = 62 % 60 = 2
	suggestion := suggestJitterMinute([]int{55})
	if suggestion == "" {
		t.Fatal("expected a non-empty suggestion")
	}
}

func TestCheckJitter_EmptyInput(t *testing.T) {
	warnings := CheckJitter([]Job{})
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings for empty input, got %d", len(warnings))
	}
}
