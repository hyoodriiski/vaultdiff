package diff

import (
	"regexp"
	"strings"
)

// RedactOptions controls key-based value redaction (full suppression vs masking).
type RedactOptions struct {
	Enabled  bool
	Patterns []*regexp.Regexp
}

// DefaultRedactOptions returns redaction disabled by default.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Enabled: false,
	}
}

const redactedPlaceholder = "[REDACTED]"

// ApplyRedact replaces entire values (old and new) for keys matching any
// redact pattern. Unlike masking, the original value is never revealed.
func ApplyRedact(changes []Change, opts RedactOptions) []Change {
	if !opts.Enabled || len(opts.Patterns) == 0 {
		return changes
	}

	result := make([]Change, len(changes))
	for i, c := range changes {
		if matchesAny(c.Key, opts.Patterns) {
			c.OldValue = redactValue(c.OldValue)
			c.NewValue = redactValue(c.NewValue)
		}
		result[i] = c
	}
	return result
}

func redactValue(v string) string {
	if v == "" {
		return ""
	}
	return redactedPlaceholder
}

func matchesAny(key string, patterns []*regexp.Regexp) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if p.MatchString(lower) {
			return true
		}
	}
	return false
}
