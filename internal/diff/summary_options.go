package diff

import (
	"fmt"
	"strings"
)

// ParseSummaryFlags parses a comma-separated list of summary option names
// and returns a populated SummaryOptions.
// Supported tokens: "counts", "percents".
func ParseSummaryFlags(flag string) (SummaryOptions, error) {
	opts := SummaryOptions{}
	if strings.TrimSpace(flag) == "" {
		return DefaultSummaryOptions(), nil
	}

	parts := strings.Split(flag, ",")
	for _, part := range parts {
		token := strings.TrimSpace(strings.ToLower(part))
		switch token {
		case "counts":
			opts.ShowCounts = true
		case "percents":
			opts.ShowPercents = true
		default:
			return SummaryOptions{}, fmt.Errorf("unknown summary option: %q (supported: counts, percents)", token)
		}
	}
	return opts, nil
}

// SupportedSummaryOptions lists tokens accepted by ParseSummaryFlags.
var SupportedSummaryOptions = []string{"counts", "percents"}
