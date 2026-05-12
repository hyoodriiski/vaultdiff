package diff

// ContextOptions controls how many surrounding unchanged lines
// are shown around changes in a unified-diff style output.
type ContextOptions struct {
	// Lines is the number of unchanged neighbours to include on each side.
	// A value of 0 means no context; negative values are treated as 0.
	Lines int

	// Enabled reports whether context output is active.
	Enabled bool
}

// DefaultContextOptions returns the default context configuration.
func DefaultContextOptions() ContextOptions {
	return ContextOptions{
		Lines:   3,
		Enabled: false,
	}
}

// ApplyContext takes a slice of Change values and, given a set of
// ContextOptions, returns a filtered slice that contains every
// non-Unchanged change plus up to opts.Lines Unchanged neighbours on
// either side of each changed entry.
//
// If opts.Enabled is false the original slice is returned unmodified.
func ApplyContext(changes []Change, opts ContextOptions) []Change {
	if !opts.Enabled || opts.Lines < 0 {
		return changes
	}

	n := len(changes)
	if n == 0 {
		return changes
	}

	// Mark indices that must be included.
	include := make([]bool, n)

	for i, c := range changes {
		if c.Type == ChangeTypeUnchanged {
			continue
		}
		// Include the change itself and its neighbours.
		start := i - opts.Lines
		if start < 0 {
			start = 0
		}
		end := i + opts.Lines
		if end >= n {
			end = n - 1
		}
		for j := start; j <= end; j++ {
			include[j] = true
		}
	}

	out := make([]Change, 0, n)
	for i, c := range changes {
		if include[i] {
			out = append(out, c)
		}
	}
	return out
}
