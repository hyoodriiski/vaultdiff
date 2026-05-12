package diff

import "strings"

// ParseAnnotateFlags parses CLI flag strings into AnnotateOptions.
// sourceFlag format: "a=path/one,b=path/two"
func ParseAnnotateFlags(enableIndex bool, enableSource bool, sourceFlag string) (AnnotateOptions, error) {
	opts := DefaultAnnotateOptions()

	if enableIndex || enableSource {
		opts.Enabled = true
	}

	opts.ShowIndex = enableIndex
	opts.ShowSource = enableSource

	if sourceFlag == "" {
		return opts, nil
	}

	parts := strings.Split(sourceFlag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		switch key {
		case "a":
			opts.SourceA = val
		case "b":
			opts.SourceB = val
		}
	}

	return opts, nil
}
