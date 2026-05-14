package analyzer

// CatchupWarning describes a job that may produce a large number of catch-up
// executions when the scheduler restarts after a prolonged outage.
type CatchupWarning struct {
	// Job is the job that triggered the warning.
	Job Job
	// FiresIn24h is the number of times the job fires in a 24-hour window.
	FiresIn24h int
	// EstimatedCatchup is the estimated number of missed runs after OutageHours
	// hours of downtime.
	EstimatedCatchup int
	// OutageHours is the assumed outage duration used for the calculation.
	OutageHours int
}

// defaultOutageHours is the assumed scheduler downtime used when the caller
// does not supply an explicit value.
const defaultOutageHours = 8

// CheckCatchup inspects each job and warns when the number of missed runs
// after a scheduler outage of outageHours hours exceeds threshold.
//
// Pass outageHours <= 0 to use the built-in default of 8 hours.
// Pass threshold <= 0 to use the built-in default of 10 missed runs.
func CheckCatchup(jobs []Job, outageHours, threshold int) []CatchupWarning {
	if outageHours <= 0 {
		outageHours = defaultOutageHours
	}
	if threshold <= 0 {
		threshold = 10
	}

	var warnings []CatchupWarning
	for _, j := range jobs {
		if j.Schedule == nil {
			continue
		}

		fires := firesIn24hWindow(j)
		if fires == 0 {
			continue
		}

		// Estimate missed runs: fires-per-minute * outage-minutes.
		outageMinutes := outageHours * 60
		missed := (fires * outageMinutes) / (24 * 60)
		if missed < 1 {
			missed = 1
		}

		if missed >= threshold {
			warnings = append(warnings, CatchupWarning{
				Job:              j,
				FiresIn24h:       fires,
				EstimatedCatchup: missed,
				OutageHours:      outageHours,
			})
		}
	}
	return warnings
}

// firesIn24hWindow counts how many distinct (hour, minute) pairs in a 24-hour
// day the job's schedule covers.
func firesIn24hWindow(j Job) int {
	count := 0
	for h := 0; h < 24; h++ {
		for _, hour := range j.Schedule.Hours {
			if hour == h {
				count += len(j.Schedule.Minutes)
				break
			}
		}
	}
	return count
}
