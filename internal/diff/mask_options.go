package diff

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseMaskPatterns parses a comma-separated list of regex patterns into MaskOptions.
// An empty string returns the default options with masking disabled.
func ParseMaskPatterns(raw string) (MaskOptions, error) {
	opts := DefaultMaskOptions()
	if raw == "" {
		return opts, nil
	}
	parts := strings.Split(raw, ",")
	var patterns []*regexp.Regexp
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return opts, fmt.Errorf("invalid mask pattern %q: %w", p, err)
		}
		patterns = append(patterns, re)
	}
	if len(patterns) > 0 {
		opts.Patterns = patterns
	}
	opts.Enabled = true
	return opts, nil
}
