package diff

import (
	"testing"
)

func makeScoreChanges() []Change {
	return []Change{
		{Key: "a", Type: ChangeTypeUnchanged},
		{Key: "b", Type: ChangeTypeUnchanged},
		{Key: "c", Type: ChangeTypeModified, OldValue: "x", NewValue: "y"},
		{Key: "d", Type: ChangeTypeAdded, NewValue: "new"},
	}
}

func TestComputeScore_Disabled(t *testing.T) {
	opts := DefaultScoreOptions()
	result := ComputeScore(makeScoreChanges(), opts)
	if result.Score != -1 {
		t.Errorf("expected Score=-1 when disabled, got %f", result.Score)
	}
}

func TestComputeScore_AllUnchanged(t *testing.T) {
	changes := []Change{
		{Key: "a", Type: ChangeTypeUnchanged},
		{Key: "b", Type: ChangeTypeUnchanged},
	}
	opts := DefaultScoreOptions()
	opts.Enabled = true
	result := ComputeScore(changes, opts)
	if result.Score != 1.0 {
		t.Errorf("expected Score=1.0 for all unchanged, got %f", result.Score)
	}
}

func TestComputeScore_AllAdded(t *testing.T) {
	changes := []Change{
		{Key: "a", Type: ChangeTypeAdded},
		{Key: "b", Type: ChangeTypeAdded},
	}
	opts := DefaultScoreOptions()
	opts.Enabled = true
	result := ComputeScore(changes, opts)
	if result.Score != 0.0 {
		t.Errorf("expected Score=0.0 for all added, got %f", result.Score)
	}
}

func TestComputeScore_ModifiedHalfWeight(t *testing.T) {
	changes := []Change{
		{Key: "a", Type: ChangeTypeUnchanged},
		{Key: "b", Type: ChangeTypeModified, OldValue: "x", NewValue: "y"},
	}
	opts := DefaultScoreOptions()
	opts.Enabled = true
	opts.ModifiedWeight = 0.5
	result := ComputeScore(changes, opts)
	// weightedChanged=0.5, total=2 → score=1-(0.5/2)=0.75
	if result.Score != 0.75 {
		t.Errorf("expected Score=0.75, got %f", result.Score)
	}
}

func TestComputeScore_Empty(t *testing.T) {
	opts := DefaultScoreOptions()
	opts.Enabled = true
	result := ComputeScore([]Change{}, opts)
	if result.Score != -1 {
		t.Errorf("expected Score=-1 for empty changes, got %f", result.Score)
	}
}

func TestParseScoreFlags_Enabled(t *testing.T) {
	opts, err := ParseScoreFlags([]string{"enabled"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestParseScoreFlags_ModifiedWeight(t *testing.T) {
	opts, err := ParseScoreFlags([]string{"enabled", "modified-weight=0.3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.ModifiedWeight != 0.3 {
		t.Errorf("expected ModifiedWeight=0.3, got %f", opts.ModifiedWeight)
	}
}

func TestParseScoreFlags_InvalidWeight(t *testing.T) {
	_, err := ParseScoreFlags([]string{"added-weight=1.5"})
	if err == nil {
		t.Error("expected error for weight > 1.0")
	}
}

func TestParseScoreFlags_UnknownFlag(t *testing.T) {
	_, err := ParseScoreFlags([]string{"unknown-flag=1"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}
