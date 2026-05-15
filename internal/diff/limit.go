package diff

import "fmt"

// LimitOptions controls how many changes are returned in the output.
type LimitOptions struct {
	Enabled  bool
	MaxItems int
	Truncated bool
}

// DefaultLimitOptions returns a LimitOptions with limiting disabled.
func DefaultLimitOptions() LimitOptions {
	return LimitOptions{
		Enabled:  false,
		MaxItems: 100,
	}
}

// ApplyLimit trims the change list to at most MaxItems entries when enabled.
// It sets Truncated on the returned options if items were dropped.
func ApplyLimit(changes []Change, opts LimitOptions) ([]Change, LimitOptions) {
	if !opts.Enabled || opts.MaxItems <= 0 {
		return changes, opts
	}

	if len(changes) <= opts.MaxItems {
		opts.Truncated = false
		return changes, opts
	}

	opts.Truncated = true
	return changes[:opts.MaxItems], opts
}

// LimitSummary returns a human-readable note when the output was truncated.
func LimitSummary(opts LimitOptions, total int) string {
	if !opts.Truncated {
		return ""
	}
	return fmt.Sprintf("output limited to %d of %d changes", opts.MaxItems, total)
}
