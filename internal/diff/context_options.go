package diff

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseContextFlag parses a string such as "3" or "5" into a
// ContextOptions value with Enabled set to true.
//
// An empty string disables context and returns DefaultContextOptions
// with Enabled=false.
//
// Returns an error when the value is not a non-negative integer.
func ParseContextFlag(s string) (ContextOptions, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		opts := DefaultContextOptions()
		opts.Enabled = false
		return opts, nil
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return ContextOptions{}, fmt.Errorf("invalid context value %q: must be a non-negative integer", s)
	}
	if v < 0 {
		return ContextOptions{}, fmt.Errorf("invalid context value %q: must be >= 0", s)
	}

	return ContextOptions{
		Lines:   v,
		Enabled: true,
	}, nil
}
