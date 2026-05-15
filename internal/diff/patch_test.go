package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func makePatchChanges() []Change {
	return []Change{
		{Key: "db_host", Type: ChangeTypeAdded, NewValue: "localhost"},
		{Key: "db_pass", Type: ChangeTypeRemoved, OldValue: "secret"},
		{Key: "db_port", Type: ChangeTypeModified, OldValue: "5432", NewValue: "5433"},
		{Key: "app_name", Type: ChangeTypeUnchanged, OldValue: "app", NewValue: "app"},
	}
}

func TestGeneratePatch_Disabled(t *testing.T) {
	opts := DefaultPatchOptions()
	patch, err := GeneratePatch(makePatchChanges(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if patch != nil {
		t.Errorf("expected nil patch when disabled, got %+v", patch)
	}
}

func TestGeneratePatch_OnlyActionableChanges(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.Enabled = true
	opts.DryRun = true

	patch, err := GeneratePatch(makePatchChanges(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if patch == nil {
		t.Fatal("expected non-nil patch")
	}
	// unchanged should be excluded
	if len(patch.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(patch.Entries))
	}
}

func TestGeneratePatch_OpsCorrect(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.Enabled = true
	opts.DryRun = true

	patch, _ := GeneratePatch(makePatchChanges(), opts)

	expected := map[string]string{
		"db_host": "set",
		"db_pass": "delete",
		"db_port": "set",
	}
	for _, e := range patch.Entries {
		if op, ok := expected[e.Key]; ok {
			if e.Op != op {
				t.Errorf("key %q: expected op %q, got %q", e.Key, op, e.Op)
			}
		}
	}
}

func TestGeneratePatch_WritesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "patch.json")

	opts := DefaultPatchOptions()
	opts.Enabled = true
	opts.OutputPath = out

	_, err := GeneratePatch(makePatchChanges(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	var p Patch
	if err := json.Unmarshal(data, &p); err != nil {
		t.Fatalf("invalid JSON in patch file: %v", err)
	}
	if len(p.Entries) == 0 {
		t.Error("expected entries in patch file")
	}
	if p.GeneratedAt == "" {
		t.Error("expected GeneratedAt to be set")
	}
}

func TestGeneratePatch_DryRunNoFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "patch.json")

	opts := DefaultPatchOptions()
	opts.Enabled = true
	opts.OutputPath = out
	opts.DryRun = true

	_, err := GeneratePatch(makePatchChanges(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("expected no file written in dry-run mode")
	}
}
