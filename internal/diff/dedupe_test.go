package diff

import (
	"testing"
)

func makeDedupeChanges() []Change {
	return []Change{
		{Key: "alpha", ChangeType: "added", OldValue: "", NewValue: "1"},
		{Key: "beta", ChangeType: "modified", OldValue: "old", NewValue: "new"},
		{Key: "alpha", ChangeType: "added", OldValue: "", NewValue: "1"},   // duplicate
		{Key: "gamma", ChangeType: "removed", OldValue: "x", NewValue: ""},
		{Key: "beta", ChangeType: "modified", OldValue: "old", NewValue: "new"}, // duplicate
		{Key: "beta", ChangeType: "modified", OldValue: "old", NewValue: "different"}, // not a duplicate — different NewValue
	}
}

func TestApplyDedupe_Disabled(t *testing.T) {
	changes := makeDedupeChanges()
	opts := DedupeOptions{Enabled: false}
	result := ApplyDedupe(changes, opts)
	if len(result) != len(changes) {
		t.Errorf("expected %d changes when disabled, got %d", len(changes), len(result))
	}
}

func TestApplyDedupe_RemovesDuplicates(t *testing.T) {
	changes := makeDedupeChanges()
	opts := DedupeOptions{Enabled: true}
	result := ApplyDedupe(changes, opts)
	// 6 input, 2 exact duplicates removed → 4 unique
	if len(result) != 4 {
		t.Errorf("expected 4 unique changes, got %d", len(result))
	}
}

func TestApplyDedupe_PreservesOrder(t *testing.T) {
	changes := []Change{
		{Key: "z", ChangeType: "added", NewValue: "1"},
		{Key: "a", ChangeType: "added", NewValue: "2"},
		{Key: "z", ChangeType: "added", NewValue: "1"}, // duplicate
	}
	opts := DedupeOptions{Enabled: true}
	result := ApplyDedupe(changes, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result))
	}
	if result[0].Key != "z" || result[1].Key != "a" {
		t.Errorf("expected order [z, a], got [%s, %s]", result[0].Key, result[1].Key)
	}
}

func TestApplyDedupe_Empty(t *testing.T) {
	opts := DedupeOptions{Enabled: true}
	result := ApplyDedupe([]Change{}, opts)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestDefaultDedupeOptions(t *testing.T) {
	opts := DefaultDedupeOptions()
	if opts.Enabled {
		t.Error("expected deduplication to be disabled by default")
	}
}
