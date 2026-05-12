package diff_test

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
)

// TestSortAfterFilter verifies Sort works correctly on filtered output.
func TestSortAfterFilter(t *testing.T) {
	changes := []diff.Change{
		{Key: "z_token", Type: diff.ChangeTypeAdded, NewValue: "abc"},
		{Key: "a_password", Type: diff.ChangeTypeRemoved, OldValue: "secret"},
		{Key: "m_host", Type: diff.ChangeTypeModified, OldValue: "old", NewValue: "new"},
		{Key: "b_port", Type: diff.ChangeTypeUnchanged, OldValue: "5432", NewValue: "5432"},
	}

	filterOpts := diff.DefaultFilterOptions()
	filterOpts.IncludeUnchanged = false
	filtered := diff.Filter(changes, filterOpts)

	sortOpts, err := diff.ParseSortFlag("key")
	if err != nil {
		t.Fatalf("ParseSortFlag error: %v", err)
	}
	sorted := diff.Sort(filtered, sortOpts)

	if len(sorted) != 3 {
		t.Fatalf("expected 3 changes after filter, got %d", len(sorted))
	}

	expectedKeys := []string{"a_password", "m_host", "z_token"}
	for i, c := range sorted {
		if c.Key != expectedKeys[i] {
			t.Errorf("index %d: got %q, want %q", i, c.Key, expectedKeys[i])
		}
	}
}

// TestSortByTypeAfterCompare verifies Sort integrates with Compare output.
func TestSortByTypeAfterCompare(t *testing.T) {
	old := map[string]interface{}{
		"removed_key": "gone",
		"shared_key":  "same",
	}
	new := map[string]interface{}{
		"added_key":  "new",
		"shared_key": "same",
	}

	changes := diff.Compare(old, new)
	sortOpts, _ := diff.ParseSortFlag("type")
	sorted := diff.Sort(changes, sortOpts)

	if len(sorted) == 0 {
		t.Fatal("expected changes from Compare")
	}

	// Verify ordering: Added before Removed before Unchanged
	for i := 1; i < len(sorted); i++ {
		prev := sorted[i-1].Type
		curr := sorted[i].Type
		if changeTypeOrderExported(prev) > changeTypeOrderExported(curr) {
			t.Errorf("sort order violated at index %d: %v > %v", i, prev, curr)
		}
	}
}

func changeTypeOrderExported(ct diff.ChangeType) int {
	switch ct {
	case diff.ChangeTypeAdded:
		return 0
	case diff.ChangeTypeRemoved:
		return 1
	case diff.ChangeTypeModified:
		return 2
	case diff.ChangeTypeUnchanged:
		return 3
	default:
		return 4
	}
}
