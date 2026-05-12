package diff

import (
	"regexp"
	"strings"
)

// MaskOptions controls which secret values are masked in output.
type MaskOptions struct {
	Enabled     bool
	Patterns    []*regexp.Regexp
	MaskString  string
}

// DefaultMaskOptions returns MaskOptions with common sensitive key patterns.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Enabled: false,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)password`),
			regexp.MustCompile(`(?i)secret`),
			regexp.MustCompile(`(?i)token`),
			regexp.MustCompile(`(?i)api[_-]?key`),
			regexp.MustCompile(`(?i)private[_-]?key`),
		},
		MaskString: "***",
	}
}

// isSensitiveKey returns true if the key matches any mask pattern.
func (o MaskOptions) isSensitiveKey(key string) bool {
	for _, p := range o.Patterns {
		if p.MatchString(key) {
			return true
		}
	}
	return false
}

// ApplyMask replaces sensitive values in a slice of Changes with the mask string.
func ApplyMask(changes []Change, opts MaskOptions) []Change {
	if !opts.Enabled {
		return changes
	}
	mask := opts.MaskString
	if mask == "" {
		mask = "***"
	}
	result := make([]Change, len(changes))
	for i, c := range changes {
		if opts.isSensitiveKey(c.Key) {
			c.OldValue = maskIfNonEmpty(c.OldValue, mask)
			c.NewValue = maskIfNonEmpty(c.NewValue, mask)
		}
		result[i] = c
	}
	return result
}

func maskIfNonEmpty(v, mask string) string {
	if strings.TrimSpace(v) == "" {
		return v
	}
	return mask
}
