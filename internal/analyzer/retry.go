package analyzer

// RetryWarning describes a job whose schedule fires so infrequently that a
// single missed execution may go unnoticed for a long time, making automatic
// retry logic especially important.
type RetryWarning struct {
	JobName    string
	Expression string
	// IntervalMinutes is the approximate gap (in minutes) between consecutive
	// firings of the schedule.
	IntervalMinutes int
	Message         string
}

// retryThresholdMinutes is the minimum gap that triggers a retry warning.
// Jobs that fire less than once per day (1440 min) are flagged.
const retryThresholdMinutes = 1440

// CheckRetry inspects each job and warns when the gap between consecutive
// firings exceeds retryThresholdMinutes, indicating that a missed run would
// not self-correct quickly.
func CheckRetry(jobs []Job) []RetryWarning {
	var warnings []RetryWarning

	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}

		interval := approximateIntervalMinutes(j.Schedule)
		if interval < retryThresholdMinutes {
			continue
		}

		warnings = append(warnings, RetryWarning{
			JobName:         j.Name,
			Expression:      j.Expression,
			IntervalMinutes: interval,
			Message: fmt.Sprintf(
				"job fires roughly every %d minutes (~%s); consider adding retry logic to avoid silent failures",
				interval, formatInterval(interval),
			),
		})
	}

	return warnings
}

// approximateIntervalMinutes returns the average gap in minutes between two
// consecutive firings derived from the expanded schedule fields.
func approximateIntervalMinutes(s *Schedule) int {
	totalSlotsPerWeek := len(s.DayOfWeek) * len(s.Hours) * len(s.Minutes)
	if totalSlotsPerWeek == 0 {
		return retryThresholdMinutes * 10 // effectively never fires
	}
	const minutesPerWeek = 7 * 24 * 60
	interval := minutesPerWeek / totalSlotsPerWeek
	if interval == 0 {
		interval = 1
	}
	return interval
}

// formatInterval converts minutes into a human-readable duration string.
func formatInterval(minutes int) string {
	switch {
	case minutes >= 10080:
		return fmt.Sprintf("%d week(s)", minutes/10080)
	case minutes >= 1440:
		return fmt.Sprintf("%d day(s)", minutes/1440)
	case minutes >= 60:
		return fmt.Sprintf("%d hour(s)", minutes/60)
	default:
		return fmt.Sprintf("%d minute(s)", minutes)
	}
}
