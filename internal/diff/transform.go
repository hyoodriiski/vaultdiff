package diff

import "strings"

// TransformOptions controls value transformation applied to changes before comparison output.
type TransformOptions struct {
	Enabled    bool
	Uppercase  bool
	Lowercase  bool
	TrimPrefix string
	TrimSuffix string
}

// DefaultTransformOptions returns a no-op TransformOptions.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{
		Enabled: false,
	}
}

// ApplyTransform applies value transformations to all Change entries.
// Transformations are applied to OldValue and NewValue for modified/unchanged changes,
// and to the relevant value for added/removed changes.
func ApplyTransform(changes []Change, opts TransformOptions) []Change {
	if !opts.Enabled {
		return changes
	}

	result := make([]Change, len(changes))
	for i, c := range changes {
		c.OldValue = transformValue(c.OldValue, opts)
		c.NewValue = transformValue(c.NewValue, opts)
		result[i] = c
	}
	return result
}

func transformValue(v string, opts TransformOptions) string {
	if v == "" {
		return v
	}
	if opts.TrimPrefix != "" {
		v = strings.TrimPrefix(v, opts.TrimPrefix)
	}
	if opts.TrimSuffix != "" {
		v = strings.TrimSuffix(v, opts.TrimSuffix)
	}
	if opts.Uppercase {
		v = strings.ToUpper(v)
	}
	if opts.Lowercase {
		v = strings.ToLower(v)
	}
	return v
}
