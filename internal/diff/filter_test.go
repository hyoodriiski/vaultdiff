package diff

import (
	"testing"
)

func makeTestChanges() []Change {
	return []Change{
		{Key: "added_key", OldValue: "", NewValue: "new", Type: string(ChangeTypeAdded)},
		{Key: "removed_key", OldValue: "old", NewValue: "", Type: string(ChangeTypeRemoved)},
		{Key: "modified_key", OldValue: "old", NewValue: "new", Type: string(ChangeTypeModified)},
		{Key: "unchanged_key", OldValue: "same", NewValue: "same", Type: string(ChangeTypeUnchanged)},
	}
}

func TestFilter_DefaultOptions(t *testing.T) {
	changes := makeTestChanges()
	opts := DefaultFilterOptions()
	result := Filter(changes, opts)

	if len(result) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(result))
	}
	for _, c := range result {
		if ChangeType(c.Type) == ChangeTypeUnchanged {
			t.Errorf("unexpected unchanged entry in default filter result")
		}
	}
}

func TestFilter_OnlyAdded(t *testing.T) {
	changes := makeTestChanges()
	opts := FilterOptions{IncludeAdded: true}
	result := Filter(changes, opts)

	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if ChangeType(result[0].Type) != ChangeTypeAdded {
		t.Errorf("expected added change, got %s", result[0].Type)
	}
}

func TestFilter_OnlyRemoved(t *testing.T) {
	changes := makeTestChanges()
	opts := FilterOptions{IncludeRemoved: true}
	result := Filter(changes, opts)

	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if ChangeType(result[0].Type) != ChangeTypeRemoved {
		t.Errorf("expected removed change, got %s", result[0].Type)
	}
}

func TestFilter_IncludeUnchanged(t *testing.T) {
	changes := makeTestChanges()
	opts := FilterOptions{IncludeUnchanged: true}
	result := Filter(changes, opts)

	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if ChangeType(result[0].Type) != ChangeTypeUnchanged {
		t.Errorf("expected unchanged change, got %s", result[0].Type)
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	result := Filter([]Change{}, DefaultFilterOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}

func TestFilter_AllOptions(t *testing.T) {
	changes := makeTestChanges()
	opts := FilterOptions{
		IncludeAdded:    true,
		IncludeRemoved:  true,
		IncludeModified: true,
		IncludeUnchanged: true,
	}
	result := Filter(changes, opts)
	if len(result) != len(changes) {
		t.Errorf("expected %d changes, got %d", len(changes), len(result))
	}
}
