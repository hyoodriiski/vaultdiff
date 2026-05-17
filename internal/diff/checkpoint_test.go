package diff

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeCheckpointChanges() []Change {
	return []Change{
		{Key: "alpha", OldValue: "1", NewValue: "2", Type: ChangeModified},
		{Key: "beta", OldValue: "", NewValue: "new", Type: ChangeAdded},
	}
}

func TestSaveCheckpoint_Disabled(t *testing.T) {
	opts := DefaultCheckpointOptions()
	if err := SaveCheckpoint(opts, makeCheckpointChanges()); err != nil {
		t.Fatalf("expected no error when disabled, got %v", err)
	}
}

func TestSaveCheckpoint_WritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "checkpoint.json")

	opts := CheckpointOptions{Enabled: true, Path: path}
	if err := SaveCheckpoint(opts, makeCheckpointChanges()); err != nil {
		t.Fatalf("SaveCheckpoint: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestLoadCheckpoint_ReturnsChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "checkpoint.json")

	opts := CheckpointOptions{Enabled: true, Path: path}
	input := makeCheckpointChanges()

	if err := SaveCheckpoint(opts, input); err != nil {
		t.Fatalf("SaveCheckpoint: %v", err)
	}

	got, savedAt, err := LoadCheckpoint(opts)
	if err != nil {
		t.Fatalf("LoadCheckpoint: %v", err)
	}
	if len(got) != len(input) {
		t.Errorf("expected %d changes, got %d", len(input), len(got))
	}
	if savedAt.IsZero() {
		t.Error("expected non-zero saved_at timestamp")
	}
	if time.Since(savedAt) > 5*time.Second {
		t.Errorf("saved_at seems too old: %v", savedAt)
	}
}

func TestLoadCheckpoint_MissingFile(t *testing.T) {
	opts := CheckpointOptions{Enabled: true, Path: "/nonexistent/path/cp.json"}
	changes, ts, err := LoadCheckpoint(opts)
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if changes != nil {
		t.Errorf("expected nil changes, got %v", changes)
	}
	if !ts.IsZero() {
		t.Errorf("expected zero timestamp, got %v", ts)
	}
}

func TestLoadCheckpoint_Disabled(t *testing.T) {
	opts := DefaultCheckpointOptions()
	changes, _, err := LoadCheckpoint(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changes != nil {
		t.Errorf("expected nil changes when disabled")
	}
}

func TestParseCheckpointFlags_Disabled(t *testing.T) {
	opts, err := ParseCheckpointFlags(false, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected disabled")
	}
}

func TestParseCheckpointFlags_MissingPath(t *testing.T) {
	_, err := ParseCheckpointFlags(true, "", false)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestParseCheckpointFlags_Valid(t *testing.T) {
	opts, err := ParseCheckpointFlags(true, "/tmp/cp.json", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected enabled")
	}
	if opts.Path != "/tmp/cp.json" {
		t.Errorf("unexpected path: %s", opts.Path)
	}
	if !opts.AutoLoad {
		t.Error("expected auto-load true")
	}
}
