package diff

import "strings"

// TagOptions controls how changes are tagged with arbitrary labels.
type TagOptions struct {
	Enabled bool
	// Tags is a map of key-pattern (substring) to label.
	// If a change's Key contains the pattern, the label is appended to Annotations.
	Tags map[string]string
}

// DefaultTagOptions returns tagging disabled.
func DefaultTagOptions() TagOptions {
	return TagOptions{
		Enabled: false,
		Tags:    map[string]string{},
	}
}

// ApplyTag iterates over changes and appends matching tag labels to each
// change's Annotations slice. Multiple tags may match a single change.
func ApplyTag(changes []Change, opts TagOptions) []Change {
	if !opts.Enabled || len(opts.Tags) == 0 {
		return changes
	}

	result := make([]Change, len(changes))
	for i, c := range changes {
		for pattern, label := range opts.Tags {
			if strings.Contains(c.Key, pattern) {
				c.Annotations = appendUnique(c.Annotations, label)
			}
		}
		result[i] = c
	}
	return result
}

// appendUnique appends s to slice only if it is not already present.
func appendUnique(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}

// ParseTagFlags parses a slice of "pattern=label" strings into a TagOptions.
func ParseTagFlags(flags []string) (TagOptions, error) {
	opts := DefaultTagOptions()
	if len(flags) == 0 {
		return opts, nil
	}
	opts.Enabled = true
	for _, flag := range flags {
		parts := strings.SplitN(flag, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			continue
		}
		opts.Tags[parts[0]] = parts[1]
	}
	return opts, nil
}
