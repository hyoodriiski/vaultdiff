package diff

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseLimitFlag parses a --limit flag value into LimitOptions.
// An empty string disables limiting. A positive integer enables it.
// Returns an error for non-positive or non-numeric values.
func ParseLimitFlag(value string) (LimitOptions, error) {
	opts := DefaultLimitOptions()

	value = strings.TrimSpace(value)
	if value == "" || value == "0" {
		opts.Enabled = false
		return opts, nil
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return opts, fmt.Errorf("invalid limit value %q: must be a non-negative integer", value)
	}

	if n < 0 {
		return opts, fmt.Errorf("invalid limit value %d: must be >= 0", n)
	}

	opts.Enabled = true
	opts.MaxItems = n
	return opts, nil
}
