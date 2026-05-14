package diff

// DedupeOptions controls deduplication behavior for change entries.
type DedupeOptions struct {
	Enabled bool
}

// DefaultDedupeOptions returns sensible defaults (deduplication disabled).
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{
		Enabled: false,
	}
}

// ApplyDedupe removes duplicate Change entries from the slice.
// Two changes are considered duplicates if they share the same Key, ChangeType,
// OldValue, and NewValue. The first occurrence is kept; subsequent ones are dropped.
func ApplyDedupe(changes []Change, opts DedupeOptions) []Change {
	if !opts.Enabled || len(changes) == 0 {
		return changes
	}

	type fingerprint struct {
		Key        string
		ChangeType string
		OldValue   string
		NewValue   string
	}

	seen := make(map[fingerprint]struct{}, len(changes))
	result := make([]Change, 0, len(changes))

	for _, c := range changes {
		fp := fingerprint{
			Key:        c.Key,
			ChangeType: c.ChangeType,
			OldValue:   c.OldValue,
			NewValue:   c.NewValue,
		}
		if _, exists := seen[fp]; !exists {
			seen[fp] = struct{}{}
			result = append(result, c)
		}
	}

	return result
}
