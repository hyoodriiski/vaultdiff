package diff

import "fmt"

// ThresholdOptions controls whether the diff run should fail (return an error)
// when the number or percentage of changes exceeds configured limits.
type ThresholdOptions struct {
	Enabled        bool
	MaxChanges     int     // 0 = unlimited
	MaxChangesPct  float64 // 0 = unlimited; percentage of total keys
}

// DefaultThresholdOptions returns a disabled ThresholdOptions.
func DefaultThresholdOptions() ThresholdOptions {
	return ThresholdOptions{Enabled: false}
}

// ThresholdResult holds the outcome of a threshold check.
type ThresholdResult struct {
	Exceeded       bool
	Reason         string
	ActualChanges  int
	ActualPct      float64
}

// CheckThreshold evaluates changes against the configured thresholds.
// totalKeys is the union of keys across both secrets (used for percentage calc).
func CheckThreshold(changes []Change, totalKeys int, opts ThresholdOptions) ThresholdResult {
	if !opts.Enabled {
		return ThresholdResult{}
	}

	actionable := 0
	for _, c := range changes {
		if c.Type != Unchanged {
			actionable++
		}
	}

	var pct float64
	if totalKeys > 0 {
		pct = float64(actionable) / float64(totalKeys) * 100.0
	}

	if opts.MaxChanges > 0 && actionable > opts.MaxChanges {
		return ThresholdResult{
			Exceeded:      true,
			Reason:        fmt.Sprintf("change count %d exceeds max %d", actionable, opts.MaxChanges),
			ActualChanges: actionable,
			ActualPct:     pct,
		}
	}

	if opts.MaxChangesPct > 0 && pct > opts.MaxChangesPct {
		return ThresholdResult{
			Exceeded:      true,
			Reason:        fmt.Sprintf("change percentage %.2f%% exceeds max %.2f%%", pct, opts.MaxChangesPct),
			ActualChanges: actionable,
			ActualPct:     pct,
		}
	}

	return ThresholdResult{ActualChanges: actionable, ActualPct: pct}
}
