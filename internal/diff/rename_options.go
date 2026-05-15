package diff

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseRenameFlags parses rename-related CLI flag values into RenameOptions.
// enabled: "true" / "false"
// similarity: float string between 0 and 1, e.g. "0.8"
func ParseRenameFlags(enabled, similarity string) (RenameOptions, error) {
	opts := DefaultRenameOptions()

	enabled = strings.TrimSpace(enabled)
	if enabled != "" {
		switch strings.ToLower(enabled) {
		case "true", "1", "yes":
			opts.Enabled = true
		case "false", "0", "no":
			opts.Enabled = false
		default:
			return opts, fmt.Errorf("invalid rename enabled value: %q", enabled)
		}
	}

	similarity = strings.TrimSpace(similarity)
	if similarity != "" {
		v, err := strconv.ParseFloat(similarity, 64)
		if err != nil {
			return opts, fmt.Errorf("invalid similarity value: %q", similarity)
		}
		if v < 0.0 || v > 1.0 {
			return opts, fmt.Errorf("similarity must be between 0.0 and 1.0, got %v", v)
		}
		opts.Similarity = v
	}

	return opts, nil
}
