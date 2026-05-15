package diff

import (
	"fmt"
	"strings"
)

// ParseHighlightFlags parses highlight-related CLI flags into HighlightOptions.
// enabled: whether highlighting is active.
// markers: optional "prefix:suffix" string, e.g. "[[:]]"
func ParseHighlightFlags(enabled bool, markers string) (HighlightOptions, error) {
	opts := DefaultHighlightOptions()
	opts.Enabled = enabled

	if markers == "" {
		return opts, nil
	}

	parts := strings.SplitN(markers, ":", 2)
	if len(parts) != 2 {
		return opts, fmt.Errorf("invalid highlight markers %q: expected \"prefix:suffix\" format", markers)
	}

	prefix := strings.TrimSpace(parts[0])
	suffix := strings.TrimSpace(parts[1])

	if prefix == "" || suffix == "" {
		return opts, fmt.Errorf("highlight prefix and suffix must be non-empty")
	}

	opts.Prefix = prefix
	opts.Suffix = suffix
	return opts, nil
}
