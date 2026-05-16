package diff

import (
	"testing"
)

func makeThresholdChanges() []Change {
	return []Change{
		{Key: "a", Type: Added, NewValue: "1"},
		{Key: "b", Type: Removed, OldValue: "2"},
		{Key: "c", Type: Modified, OldValue: "x", NewValue: "y"},
		{Key: "d", Type: Unchanged, OldValue: "z", NewValue: "z"},
	}
}

func TestCheckThreshold_Disabled(t *testing.T) {
	result := CheckThreshold(makeThresholdChanges(), 4, DefaultThresholdOptions())
	if result.Exceeded {
		t.Errorf("expected not exceeded when disabled")
	}
}

func TestCheckThreshold_BelowMaxChanges(t *testing.T) {
	opts := ThresholdOptions{Enabled: true, MaxChanges: 5}
	result := CheckThreshold(makeThresholdChanges(), 4, opts)
	if result.Exceeded {
		t.Errorf("expected not exceeded: got %s", result.Reason)
	}
	if result.ActualChanges != 3 {
		t.Errorf("expected 3 actionable changes, got %d", result.ActualChanges)
	}
}

func TestCheckThreshold_ExceedsMaxChanges(t *testing.T) {
	opts := ThresholdOptions{Enabled: true, MaxChanges: 2}
	result := CheckThreshold(makeThresholdChanges(), 4, opts)
	if !result.Exceeded {
		t.Errorf("expected threshold exceeded")
	}
	if result.Reason == "" {
		t.Errorf("expected non-empty reason")
	}
}

func TestCheckThreshold_ExceedsMaxPct(t *testing.T) {
	opts := ThresholdOptions{Enabled: true, MaxChangesPct: 50.0}
	// 3 actionable out of 4 total = 75%
	result := CheckThreshold(makeThresholdChanges(), 4, opts)
	if !result.Exceeded {
		t.Errorf("expected threshold exceeded at 75%%")
	}
}

func TestCheckThreshold_BelowMaxPct(t *testing.T) {
	opts := ThresholdOptions{Enabled: true, MaxChangesPct: 80.0}
	result := CheckThreshold(makeThresholdChanges(), 4, opts)
	if result.Exceeded {
		t.Errorf("expected not exceeded: got %s", result.Reason)
	}
}

func TestParseThresholdFlags_Empty(t *testing.T) {
	opts, err := ParseThresholdFlags("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Errorf("expected disabled when no flags provided")
	}
}

func TestParseThresholdFlags_MaxChanges(t *testing.T) {
	opts, err := ParseThresholdFlags("10", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled || opts.MaxChanges != 10 {
		t.Errorf("expected enabled with MaxChanges=10, got %+v", opts)
	}
}

func TestParseThresholdFlags_MaxPctWithSuffix(t *testing.T) {
	opts, err := ParseThresholdFlags("", "25%")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.MaxChangesPct != 25.0 {
		t.Errorf("expected MaxChangesPct=25.0, got %f", opts.MaxChangesPct)
	}
}

func TestParseThresholdFlags_InvalidMaxChanges(t *testing.T) {
	_, err := ParseThresholdFlags("abc", "")
	if err == nil {
		t.Errorf("expected error for invalid max-changes")
	}
}

func TestParseThresholdFlags_InvalidMaxPct(t *testing.T) {
	_, err := ParseThresholdFlags("", "150")
	if err == nil {
		t.Errorf("expected error for out-of-range max-changes-pct")
	}
}
