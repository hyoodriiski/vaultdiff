package diff

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseRedactPatterns parses a comma-separated list of regex patterns used for
// full value redaction. Returns an error if any pattern is invalid.
func ParseRedactPatterns(raw string) ([]string, []*regexp.Regexp, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil, nil
	}

	parts := strings.Split(raw, ",")
	patterns := make([]string, 0, len(parts))
	compiled := make([]*regexp.Regexp, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		re, err := regexp.Compile("(?i)" + p)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid redact pattern %q: %w", p, err)
		}
		patterns = append(patterns, p)
		compiled = append(compiled, re)
	}

	return patterns, compiled, nil
}

// SupportedRedactPatterns lists commonly used default redact key patterns.
var SupportedRedactPatterns = []string{
	"private_key",
	"secret",
	"token",
	"credential",
}
