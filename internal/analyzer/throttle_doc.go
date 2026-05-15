// Package analyzer – throttle module
//
// CheckThrottle scans all jobs and flags 5-minute windows in which the
// combined number of firings across every job exceeds a configurable
// threshold.  This catches situations where many lightweight jobs are
// inadvertently scheduled to run at the same time, creating micro-bursts
// that can overwhelm downstream services or trigger rate-limiting.
//
// # Algorithm
//
// Each (hour, 5-min-block) pair is treated as an independent slot.  For
// every job whose schedule fires inside that slot the counter is
// incremented.  When the counter exceeds MaxFiringsPer5Min a
// ThrottleWarning is emitted for that slot.
//
// # Configuration
//
// Use DefaultThrottleOptions for sensible defaults or supply a custom
// ThrottleOptions to tune the threshold for your environment.
package analyzer
