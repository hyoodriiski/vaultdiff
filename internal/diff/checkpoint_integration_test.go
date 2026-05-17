package diff

import (
	"path/filepath"
	"testing"
)

// TestCheckpointRoundTrip saves then reloads changes and verifies key/type fidelity.
func TestCheckpointRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rt.json")
	opts := CheckpointOptions{Enabled: true, Path: path}

	original := []Change{
		{Key: "db_pass", OldValue: "old", NewValue: "new", Type: ChangeModified},
		{Key: "api_key", OldValue: "", NewValue: "xyz", Type: ChangeAdded},
		{Key: "legacy", OldValue: "v", NewValue: "", Type: ChangeRemoved},
	}

	if err := SaveCheckpoint(opts, original); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, _, err := LoadCheckpoint(opts)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if len(loaded) != len(original) {
		t.Fatalf("length mismatch: want %d, got %d", len(original), len(loaded))
	}

	for i, want := range original {
		got := loaded[i]
		if got.Key != want.Key {
			t.Errorf("[%d] key: want %q, got %q", i, want.Key, got.Key)
		}
		if got.Type != want.Type {
			t.Errorf("[%d] type: want %q, got %q", i, want.Type, got.Type)
		}
		if got.OldValue != want.OldValue {
			t.Errorf("[%d] old_value: want %q, got %q", i, want.OldValue, got.OldValue)
		}
		if got.NewValue != want.NewValue {
			t.Errorf("[%d] new_value: want %q, got %q", i, want.NewValue, got.NewValue)
		}
	}
}

// TestCheckpointAfterCompare exercises the full Compare → SaveCheckpoint → LoadCheckpoint path.
func TestCheckpointAfterCompare(t *testing.T) {
	path := filepath.Join(t.TempDir(), "after_compare.json")
	opts := CheckpointOptions{Enabled: true, Path: path}

	src := map[string]interface{}{"x": "1", "y": "2"}
	dst := map[string]interface{}{"x": "99", "z": "3"}

	changes := Compare(src, dst)

	if err := SaveCheckpoint(opts, changes); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, _, err := LoadCheckpoint(opts)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if len(loaded) != len(changes) {
		t.Errorf("want %d changes, got %d", len(changes), len(loaded))
	}
}
