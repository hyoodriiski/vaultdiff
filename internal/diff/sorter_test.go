package diff

import (
	"testing"
)

func makeUnsortedChanges() []Change {
	return []Change{
		{Key: "zebra", Type: ChangeTypeUnchanged},
		{Key: "apple", Type: ChangeTypeAdded},
		{Key: "mango", Type: ChangeTypeRemoved},
		{Key: "banana", Type: ChangeTypeModified},
		{Key: "cherry", Type: ChangeTypeAdded},
	}
}

func TestSort_ByKey(t *testing.T) {
	changes := makeUnsortedChanges()
	result := Sort(changes, SortOptions{Order: SortByKey})

	expected := []string{"apple", "banana", "cherry", "mango", "zebra"}
	for i, c := range result {
		if c.Key != expected[i] {
			t.Errorf("index %d: got key %q, want %q", i, c.Key, expected[i])
		}
	}
}

func TestSort_ByKeyDesc(t *testing.T) {
	changes := makeUnsortedChanges()
	result := Sort(changes, SortOptions{Order: SortByKeyDesc})

	expected := []string{"zebra", "mango", "cherry", "banana", "apple"}
	for i, c := range result {
		if c.Key != expected[i] {
			t.Errorf("index %d: got key %q, want %q", i, c.Key, expected[i])
		}
	}
}

func TestSort_ByChangeType(t *testing.T) {
	changes := makeUnsortedChanges()
	result := Sort(changes, SortOptions{Order: SortByChangeType, StableKey: true})

	// First two should be Added (sorted by key), then Removed, then Modified, then Unchanged
	if result[0].Type != ChangeTypeAdded || result[0].Key != "apple" {
		t.Errorf("expected first Added/apple, got %v/%v", result[0].Type, result[0].Key)
	}
	if result[1].Type != ChangeTypeAdded || result[1].Key != "cherry" {
		t.Errorf("expected second Added/cherry, got %v/%v", result[1].Type, result[1].Key)
	}
	if result[2].Type != ChangeTypeRemoved {
		t.Errorf("expected third Removed, got %v", result[2].Type)
	}
	if result[3].Type != ChangeTypeModified {
		t.Errorf("expected fourth Modified, got %v", result[3].Type)
	}
	if result[4].Type != ChangeTypeUnchanged {
		t.Errorf("expected fifth Unchanged, got %v", result[4].Type)
	}
}

func TestSort_Empty(t *testing.T) {
	result := Sort([]Change{}, DefaultSortOptions())
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(result))
	}
}

func TestSort_DoesNotMutateInput(t *testing.T) {
	original := makeUnsortedChanges()
	firstKey := original[0].Key
	Sort(original, DefaultSortOptions())
	if original[0].Key != firstKey {
		t.Errorf("Sort mutated the input slice")
	}
}

func TestParseSortFlag_Valid(t *testing.T) {
	cases := []struct {
		input string
		order SortOrder
	}{
		{"key", SortByKey},
		{"key-desc", SortByKeyDesc},
		{"type", SortByChangeType},
		{"", SortByKey},
	}
	for _, tc := range cases {
		opts, err := ParseSortFlag(tc.input)
		if err != nil {
			t.Errorf("ParseSortFlag(%q) unexpected error: %v", tc.input, err)
		}
		if opts.Order != tc.order {
			t.Errorf("ParseSortFlag(%q): got order %v, want %v", tc.input, opts.Order, tc.order)
		}
	}
}

func TestParseSortFlag_Invalid(t *testing.T) {
	_, err := ParseSortFlag("random")
	if err == nil {
		t.Error("expected error for invalid sort flag, got nil")
	}
}
