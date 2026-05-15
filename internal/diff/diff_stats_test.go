package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeStatsChanges() []Change {
	return []Change{
		{Key: "a", Type: ChangeTypeAdded, NewValue: "1"},
		{Key: "b", Type: ChangeTypeRemoved, OldValue: "2"},
		{Key: "c", Type: ChangeTypeModified, OldValue: "old", NewValue: "new"},
		{Key: "d", Type: ChangeTypeUnchanged, OldValue: "same", NewValue: "same"},
		{Key: "e", Type: ChangeTypeUnchanged, OldValue: "x", NewValue: "x"},
	}
}

func TestComputeStats_Disabled(t *testing.T) {
	opts := DefaultStatsOptions()
	opts.Enabled = false
	stats := ComputeStats(makeStatsChanges(), opts)
	assert.Equal(t, DiffStats{}, stats)
}

func TestComputeStats_Empty(t *testing.T) {
	stats := ComputeStats([]Change{}, DefaultStatsOptions())
	assert.Equal(t, DiffStats{}, stats)
}

func TestComputeStats_Counts(t *testing.T) {
	stats := ComputeStats(makeStatsChanges(), DefaultStatsOptions())
	assert.Equal(t, 5, stats.Total)
	assert.Equal(t, 1, stats.Added)
	assert.Equal(t, 1, stats.Removed)
	assert.Equal(t, 1, stats.Modified)
	assert.Equal(t, 2, stats.Unchanged)
}

func TestComputeStats_Percent(t *testing.T) {
	stats := ComputeStats(makeStatsChanges(), DefaultStatsOptions())
	// 3 changed out of 5 total = 60%
	assert.InDelta(t, 60.0, stats.Percent, 0.01)
}

func TestComputeStats_Ratio(t *testing.T) {
	opts := DefaultStatsOptions()
	opts.ShowRatio = true
	stats := ComputeStats(makeStatsChanges(), opts)
	assert.Equal(t, "3/5", stats.Ratio)
}

func TestComputeStats_NoRatioByDefault(t *testing.T) {
	stats := ComputeStats(makeStatsChanges(), DefaultStatsOptions())
	assert.Empty(t, stats.Ratio)
}

func TestComputeStats_AllUnchanged(t *testing.T) {
	changes := []Change{
		{Key: "a", Type: ChangeTypeUnchanged},
		{Key: "b", Type: ChangeTypeUnchanged},
	}
	stats := ComputeStats(changes, DefaultStatsOptions())
	assert.Equal(t, 0, stats.Added+stats.Removed+stats.Modified)
	assert.InDelta(t, 0.0, stats.Percent, 0.01)
}
