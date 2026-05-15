package diff

import (
	"math"
)

// ScoreOptions configures how a diff similarity score is computed.
type ScoreOptions struct {
	Enabled bool
	// Weights for each change type (0.0–1.0)
	AddedWeight    float64
	RemovedWeight  float64
	ModifiedWeight float64
}

// DefaultScoreOptions returns sensible defaults for scoring.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{
		Enabled:        false,
		AddedWeight:    1.0,
		RemovedWeight:  1.0,
		ModifiedWeight: 0.5,
	}
}

// DiffScore holds the computed similarity score for a diff result.
type DiffScore struct {
	// Score is a value between 0.0 (completely different) and 1.0 (identical).
	Score     float64
	Total     int
	Unchanged int
	Changed   int
}

// ComputeScore calculates a similarity score from a slice of Changes.
// Returns a DiffScore; if options.Enabled is false, Score is -1.
func ComputeScore(changes []Change, opts ScoreOptions) DiffScore {
	if !opts.Enabled || len(changes) == 0 {
		return DiffScore{Score: -1}
	}

	var weightedChanged float64
	unchanged := 0

	for _, c := range changes {
		switch c.Type {
		case ChangeTypeUnchanged:
			unchanged++
		case ChangeTypeAdded:
			weightedChanged += opts.AddedWeight
		case ChangeTypeRemoved:
			weightedChanged += opts.RemovedWeight
		case ChangeTypeModified:
			weightedChanged += opts.ModifiedWeight
		}
	}

	total := len(changes)
	effectiveChanged := math.Min(weightedChanged, float64(total))
	score := 1.0 - (effectiveChanged / float64(total))

	return DiffScore{
		Score:     math.Round(score*1000) / 1000,
		Total:     total,
		Unchanged: unchanged,
		Changed:   int(math.Round(effectiveChanged)),
	}
}
