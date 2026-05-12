package diff

import "fmt"

// ParseFilterFlag parses a comma-separated list of change type strings into
// a FilterOptions struct. Valid values: "added", "removed", "modified", "unchanged".
// Returns an error if an unrecognized type is provided.
func ParseFilterFlag(types []string) (FilterOptions, error) {
	if len(types) == 0 {
		return DefaultFilterOptions(), nil
	}

	opts := FilterOptions{}
	for _, t := range types {
		switch ChangeType(t) {
		case ChangeTypeAdded:
			opts.IncludeAdded = true
		case ChangeTypeRemoved:
			opts.IncludeRemoved = true
		case ChangeTypeModified:
			opts.IncludeModified = true
		case ChangeTypeUnchanged:
			opts.IncludeUnchanged = true
		default:
			return FilterOptions{}, fmt.Errorf("unknown change type filter: %q (valid: added, removed, modified, unchanged)", t)
		}
	}
	return opts, nil
}

// SupportedFilterTypes returns the list of valid filter type strings.
func SupportedFilterTypes() []string {
	return []string{
		string(ChangeTypeAdded),
		string(ChangeTypeRemoved),
		string(ChangeTypeModified),
		string(ChangeTypeUnchanged),
	}
}
