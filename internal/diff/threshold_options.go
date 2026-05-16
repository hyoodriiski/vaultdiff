package diff

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseThresholdFlags parses CLI flag strings into ThresholdOptions.
//
// maxChanges: integer string, e.g. "10"
// maxChangesPct: float string with optional % suffix, e.g. "25" or "25%"
func ParseThresholdFlags(maxChanges, maxChangesPct string) (ThresholdOptions, error) {
	opts := DefaultThresholdOptions()

	if maxChanges == "" && maxChangesPct == "" {
		return opts, nil
	}

	opts.Enabled = true

	if maxChanges != "" {
		v, err := strconv.Atoi(strings.TrimSpace(maxChanges))
		if err != nil || v < 0 {
			return opts, fmt.Errorf("invalid --max-changes value %q: must be a non-negative integer", maxChanges)
		}
		opts.MaxChanges = v
	}

	if maxChangesPct != "" {
		raw := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(maxChangesPct), "%"))
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil || v < 0 || v > 100 {
			return opts, fmt.Errorf("invalid --max-changes-pct value %q: must be a float between 0 and 100", maxChangesPct)
		}
		opts.MaxChangesPct = v
	}

	return opts, nil
}
