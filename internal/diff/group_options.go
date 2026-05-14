package diff

import (
	"fmt"
	"strings"
)

// SupportedGroupByValues lists valid group-by modes.
var SupportedGroupByValues = []string{"type", "prefix"}

// ParseGroupFlags parses CLI flag strings into GroupOptions.
// groupBy: "type" or "prefix"
// separator: string used to split key prefixes (e.g. ":" or "/")
func ParseGroupFlags(enabled bool, groupBy, separator string) (GroupOptions, error) {
	opts := DefaultGroupOptions()
	opts.Enabled = enabled

	if !enabled {
		return opts, nil
	}

	groupBy = strings.TrimSpace(strings.ToLower(groupBy))
	if groupBy == "" {
		groupBy = "type"
	}

	valid := false
	for _, v := range SupportedGroupByValues {
		if v == groupBy {
			valid = true
			break
		}
	}
	if !valid {
		return opts, fmt.Errorf("unsupported group-by value %q: must be one of %v", groupBy, SupportedGroupByValues)
	}

	opts.GroupBy = groupBy

	if separator != "" {
		opts.Separator = separator
	}

	return opts, nil
}
