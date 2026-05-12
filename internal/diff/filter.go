package diff

// ChangeType represents the type of change for a secret key.
type ChangeType string

const (
	ChangeTypeAdded    ChangeType = "added"
	ChangeTypeRemoved  ChangeType = "removed"
	ChangeTypeModified ChangeType = "modified"
	ChangeTypeUnchanged ChangeType = "unchanged"
)

// FilterOptions controls which change types are included in filtered results.
type FilterOptions struct {
	IncludeAdded    bool
	IncludeRemoved  bool
	IncludeModified bool
	IncludeUnchanged bool
}

// DefaultFilterOptions returns a FilterOptions that includes all change types
// except unchanged entries.
func DefaultFilterOptions() FilterOptions {
	return FilterOptions{
		IncludeAdded:    true,
		IncludeRemoved:  true,
		IncludeModified: true,
		IncludeUnchanged: false,
	}
}

// Filter returns a subset of changes based on the provided FilterOptions.
func Filter(changes []Change, opts FilterOptions) []Change {
	result := make([]Change, 0, len(changes))
	for _, c := range changes {
		switch ChangeType(c.Type) {
		case ChangeTypeAdded:
			if opts.IncludeAdded {
				result = append(result, c)
			}
		case ChangeTypeRemoved:
			if opts.IncludeRemoved {
				result = append(result, c)
			}
		case ChangeTypeModified:
			if opts.IncludeModified {
				result = append(result, c)
			}
		case ChangeTypeUnchanged:
			if opts.IncludeUnchanged {
				result = append(result, c)
			}
		}
	}
	return result
}
