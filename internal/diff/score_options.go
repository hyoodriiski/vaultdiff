package diff

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseScoreFlags parses CLI flag strings into ScoreOptions.
// flags example: "enabled", "modified-weight=0.3"
func ParseScoreFlags(flags []string) (ScoreOptions, error) {
	opts := DefaultScoreOptions()

	for _, f := range flags {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}

		if f == "enabled" {
			opts.Enabled = true
			continue
		}

		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 {
			return opts, fmt.Errorf("invalid score flag: %q", f)
		}

		key := strings.TrimSpace(parts[0])
		val, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return opts, fmt.Errorf("invalid weight value for %q: %w", key, err)
		}
		if val < 0.0 || val > 1.0 {
			return opts, fmt.Errorf("weight for %q must be between 0.0 and 1.0", key)
		}

		switch key {
		case "added-weight":
			opts.AddedWeight = val
		case "removed-weight":
			opts.RemovedWeight = val
		case "modified-weight":
			opts.ModifiedWeight = val
		default:
			return opts, fmt.Errorf("unknown score flag: %q", key)
		}
	}

	return opts, nil
}
