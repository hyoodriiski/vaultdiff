package diff

import "fmt"

// TruncateOptions controls value truncation in diff output.
type TruncateOptions struct {
	Enabled   bool
	MaxLength int
	Suffix    string
}

// DefaultTruncateOptions returns sensible defaults: truncation enabled at 80 chars.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		Enabled:   true,
		MaxLength: 80,
		Suffix:    "...",
	}
}

// ApplyTruncate truncates OldValue and NewValue fields in each Change if they
// exceed MaxLength. Unchanged changes are skipped when truncation is disabled.
func ApplyTruncate(changes []Change, opts TruncateOptions) []Change {
	if !opts.Enabled || opts.MaxLength <= 0 {
		return changes
	}

	result := make([]Change, len(changes))
	for i, c := range changes {
		c.OldValue = truncateString(c.OldValue, opts.MaxLength, opts.Suffix)
		c.NewValue = truncateString(c.NewValue, opts.MaxLength, opts.Suffix)
		result[i] = c
	}
	return result
}

func truncateString(s string, maxLen int, suffix string) string {
	if len(s) <= maxLen {
		return s
	}
	cutAt := maxLen - len(suffix)
	if cutAt < 0 {
		cutAt = 0
	}
	return fmt.Sprintf("%s%s", s[:cutAt], suffix)
}
