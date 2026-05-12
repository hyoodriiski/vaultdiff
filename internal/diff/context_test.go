package diff

import (
	"testing"
)

// helpers ----------------------------------------------------------------

func makeContextChanges() []Change {
	return []Change{
		{Key: "a", Type: ChangeTypeUnchanged},
		{Key: "b", Type: ChangeTypeUnchanged},
		{Key: "c", Type: ChangeTypeAdded},
		{Key: "d", Type: ChangeTypeUnchanged},
		{Key: "e", Type: ChangeTypeUnchanged},
		{Key: "f", Type: ChangeTypeUnchanged},
		{Key: "g", Type: ChangeTypeModified},
		{Key: "h", Type: ChangeTypeUnchanged},
	}
}

// ApplyContext -----------------------------------------------------------

func TestApplyContext_Disabled(t *testing.T) {
	changes := makeContextChanges()
	opts := ContextOptions{Lines: 1, Enabled: false}
	got := ApplyContext(changes, opts)
	if len(got) != len(changes) {
		t.Fatalf("expected %d changes, got %d", len(changes), len(got))
	}
}

func TestApplyContext_ZeroLines(t *testing.T) {
	changes := makeContextChanges()
	opts := ContextOptions{Lines: 0, Enabled: true}
	got := ApplyContext(changes, opts)
	// Only non-unchanged entries: c, g
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	if got[0].Key != "c" || got[1].Key != "g" {
		t.Errorf("unexpected keys: %v", got)
	}
}

func TestApplyContext_OneLinePadding(t *testing.T) {
	changes := makeContextChanges()
	opts := ContextOptions{Lines: 1, Enabled: true}
	got := ApplyContext(changes, opts)
	// Around c (index 2): b,c,d  Around g (index 6): f,g,h
	expectedKeys := []string{"b", "c", "d", "f", "g", "h"}
	if len(got) != len(expectedKeys) {
		t.Fatalf("expected %d, got %d", len(expectedKeys), len(got))
	}
	for i, k := range expectedKeys {
		if got[i].Key != k {
			t.Errorf("index %d: want %q got %q", i, k, got[i].Key)
		}
	}
}

func TestApplyContext_Empty(t *testing.T) {
	got := ApplyContext([]Change{}, ContextOptions{Lines: 3, Enabled: true})
	if len(got) != 0 {
		t.Fatalf("expected empty slice")
	}
}

// ParseContextFlag -------------------------------------------------------

func TestParseContextFlag_Empty(t *testing.T) {
	opts, err := ParseContextFlag("")
	if err != nil {
		t.Fatal(err)
	}
	if opts.Enabled {
		t.Error("expected Enabled=false for empty flag")
	}
}

func TestParseContextFlag_Valid(t *testing.T) {
	opts, err := ParseContextFlag("5")
	if err != nil {
		t.Fatal(err)
	}
	if !opts.Enabled || opts.Lines != 5 {
		t.Errorf("unexpected opts: %+v", opts)
	}
}

func TestParseContextFlag_Invalid(t *testing.T) {
	_, err := ParseContextFlag("abc")
	if err == nil {
		t.Error("expected error for non-integer")
	}
}

func TestParseContextFlag_Negative(t *testing.T) {
	_, err := ParseContextFlag("-1")
	if err == nil {
		t.Error("expected error for negative value")
	}
}
