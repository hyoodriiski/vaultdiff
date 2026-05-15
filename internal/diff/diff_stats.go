package diff

import "fmt"

// StatsOptions controls how diff statistics are computed and displayed.
type StatsOptions struct {
	Enabled     bool
	ShowPercent bool
	ShowRatio   bool
}

// DefaultStatsOptions returns sensible defaults for StatsOptions.
func DefaultStatsOptions() StatsOptions {
	return StatsOptions{
		Enabled:     true,
		ShowPercent: true,
		ShowRatio:   false,
	}
}

// DiffStats holds aggregate statistics about a set of Changes.
type DiffStats struct {
	Total     int
	Added     int
	Removed   int
	Modified  int
	Unchanged int
	Percent   float64 // percentage of keys that changed
	Ratio     string  // e.g. "3/10"
}

// ComputeStats calculates statistics over the given changes.
func ComputeStats(changes []Change, opts StatsOptions) DiffStats {
	if !opts.Enabled || len(changes) == 0 {
		return DiffStats{}
	}

	var stats DiffStats
	stats.Total = len(changes)

	for _, c := range changes {
		switch c.Type {
		case ChangeTypeAdded:
			stats.Added++
		case ChangeTypeRemoved:
			stats.Removed++
		case ChangeTypeModified:
			stats.Modified++
		case ChangeTypeUnchanged:
			stats.Unchanged++
		}
	}

	changed := stats.Added + stats.Removed + stats.Modified
	if opts.ShowPercent && stats.Total > 0 {
		stats.Percent = float64(changed) / float64(stats.Total) * 100.0
	}

	if opts.ShowRatio {
		stats.Ratio = fmt.Sprintf("%d/%d", changed, stats.Total)
	}

	return stats
}
