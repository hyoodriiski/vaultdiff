package diff

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff/types"
)

func makeNormalizeChanges() []types.Change {
	return []types.Change{
		{Key: "host", OldValue: "  localhost  ", NewValue: "localhost", ChangeType: types.ChangeTypeModified},
		{Key: "port", OldValue: "8080", NewValue: "8080", ChangeType: types.ChangeTypeUnchanged},
		{Key: "token", OldValue: "ABC", NewValue: "abc", ChangeType: types.ChangeTypeModified},
		{Key: "user", OldValue: "", NewValue: "admin", ChangeType: types.ChangeTypeAdded},
	}
}

func TestApplyNormalize_Disabled(t *testing.T) {
	changes := makeNormalizeChanges()
	opts := DefaultNormalizeOptions()
	opts.Enabled = false

	result := ApplyNormalize(changes, opts)
	if len(result) != len(changes) {
		t.Fatalf("expected %d changes, got %d", len(changes), len(result))
	}
	// values must be untouched
	if result[0].OldValue != "  localhost  " {
		t.Errorf("expected untouched OldValue, got %q", result[0].OldValue)
	}
}

func TestApplyNormalize_TrimSpace(t *testing.T) {
	changes := makeNormalizeChanges()
	opts := DefaultNormalizeOptions()
	opts.Enabled = true
	opts.TrimSpace = true

	result := ApplyNormalize(changes, opts)

	// "  localhost  " trimmed == "localhost" == NewValue → should become Unchanged
	if result[0].ChangeType != types.ChangeTypeUnchanged {
		t.Errorf("expected Unchanged after trim, got %v", result[0].ChangeType)
	}
	if result[0].OldValue != "localhost" {
		t.Errorf("expected trimmed OldValue 'localhost', got %q", result[0].OldValue)
	}
}

func TestApplyNormalize_Lowercase(t *testing.T) {
	changes := makeNormalizeChanges()
	opts := DefaultNormalizeOptions()
	opts.Enabled = true
	opts.TrimSpace = false
	opts.Lowercase = true

	result := ApplyNormalize(changes, opts)

	// "ABC" lowercased == "abc" == NewValue → should become Unchanged
	if result[2].ChangeType != types.ChangeTypeUnchanged {
		t.Errorf("expected Unchanged after lowercase, got %v", result[2].ChangeType)
	}
}

func TestApplyNormalize_AddedUnaffected(t *testing.T) {
	changes := makeNormalizeChanges()
	opts := DefaultNormalizeOptions()
	opts.Enabled = true

	result := ApplyNormalize(changes, opts)

	// Added change should remain Added regardless
	if result[3].ChangeType != types.ChangeTypeAdded {
		t.Errorf("expected Added to remain Added, got %v", result[3].ChangeType)
	}
}

func TestApplyNormalize_EmptyInput(t *testing.T) {
	opts := DefaultNormalizeOptions()
	opts.Enabled = true

	result := ApplyNormalize([]types.Change{}, opts)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}
