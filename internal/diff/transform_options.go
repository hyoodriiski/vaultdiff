package diff

import "strings"

// ParseTransformFlags parses CLI flag strings into a TransformOptions struct.
// flagStr format: "uppercase", "lowercase", "trim-prefix=foo", "trim-suffix=bar"
// Multiple flags can be passed as separate entries.
func ParseTransformFlags(flags []string) TransformOptions {
	opts := DefaultTransformOptions()
	if len(flags) == 0 {
		return opts
	}

	for _, f := range flags {
		f = strings.TrimSpace(f)
		switch {
		case f == "uppercase":
			opts.Enabled = true
			opts.Uppercase = true
		case f == "lowercase":
			opts.Enabled = true
			opts.Lowercase = true
		case strings.HasPrefix(f, "trim-prefix="):
			opts.Enabled = true
			opts.TrimPrefix = strings.TrimPrefix(f, "trim-prefix=")
		case strings.HasPrefix(f, "trim-suffix="):
			opts.Enabled = true
			opts.TrimSuffix = strings.TrimPrefix(f, "trim-suffix=")
		}
	}
	return opts
}

// SupportedTransformFlags returns the list of valid transform flag names.
func SupportedTransformFlags() []string {
	return []string{
		"uppercase",
		"lowercase",
		"trim-prefix=<value>",
		"trim-suffix=<value>",
	}
}
