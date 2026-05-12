package diff

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseTruncateFlag parses a truncate flag string of the form "<maxLen>" or
// "off" / "false" to disable truncation. Returns an error on invalid input.
//
// Examples:
//
//	"80"   -> TruncateOptions{Enabled: true, MaxLength: 80, Suffix: "..."}
//	"off"  -> TruncateOptions{Enabled: false}
//	"0"    -> TruncateOptions{Enabled: false}
func ParseTruncateFlag(flag string) (TruncateOptions, error) {
	defaults := DefaultTruncateOptions()

	trimmed := strings.TrimSpace(strings.ToLower(flag))
	if trimmed == "" {
		return defaults, nil
	}

	if trimmed == "off" || trimmed == "false" || trimmed == "0" {
		return TruncateOptions{Enabled: false, MaxLength: 0, Suffix: defaults.Suffix}, nil
	}

	n, err := strconv.Atoi(trimmed)
	if err != nil {
		return TruncateOptions{}, fmt.Errorf("invalid truncate value %q: must be a positive integer or \"off\"", flag)
	}
	if n < 0 {
		return TruncateOptions{}, fmt.Errorf("invalid truncate value %d: must be non-negative", n)
	}

	return TruncateOptions{
		Enabled:   n > 0,
		MaxLength: n,
		Suffix:    defaults.Suffix,
	}, nil
}
