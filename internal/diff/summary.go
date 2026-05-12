package diff

import "github.com/yourusername/vaultdiff/internal/diff"

// SummaryOptions controls how the summary is generated.
type SummaryOptions struct {
	ShowCounts   bool
	ShowPercents bool
}

// DefaultSummaryOptions returns sensible defaults.
func DefaultSummaryOptions() SummaryOptions {
	return SummaryOptions{
		ShowCounts:   true,
		ShowPercents: false,
	}
}

// Summary holds aggregated statistics about a set of changes.
type Summary struct {
	Total     int
	Added     int
	Removed   int
	Modified  int
	Unchanged int
}

// Percent returns the percentage of the given count relative to Total.
// Returns 0 if Total is zero.
func (s Summary) Percent(count int) float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(count) / float64(s.Total) * 100
}

// Summarize computes a Summary from a slice of Change values.
func Summarize(changes []diff.Change) Summary {
	s := Summary{Total: len(changes)}
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Modified:
			s.Modified++
		case diff.Unchanged:
			s.Unchanged++
		}
	}
	return s
}
