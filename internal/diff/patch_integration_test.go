package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPatchAfterCompare(t *testing.T) {
	src := map[string]interface{}{
		"host": "old-host",
		"port": "5432",
		"user": "admin",
	}
	dst := map[string]interface{}{
		"host": "new-host",
		"port": "5432",
		"pass": "s3cret",
	}

	changes := Compare(src, dst)

	opts := DefaultPatchOptions()
	opts.Enabled = true
	opts.DryRun = true

	patch, err := GeneratePatch(changes, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if patch == nil {
		t.Fatal("expected non-nil patch")
	}

	ops := map[string]string{}
	for _, e := range patch.Entries {
		ops[e.Key] = e.Op
	}

	if ops["host"] != "set" {
		t.Errorf("expected 'host' op=set, got %q", ops["host"])
	}
	if ops["user"] != "delete" {
		t.Errorf("expected 'user' op=delete, got %q", ops["user"])
	}
	if ops["pass"] != "set" {
		t.Errorf("expected 'pass' op=set, got %q", ops["pass"])
	}
	if _, exists := ops["port"]; exists {
		t.Error("unchanged key 'port' should not appear in patch")
	}
}

func TestPatchWriteAndReload(t *testing.T) {
	changes := []Change{
		{Key: "token", Type: ChangeTypeAdded, NewValue: "abc123"},
		{Key: "secret", Type: ChangeTypeRemoved, OldValue: "xyz"},
	}

	dir := t.TempDir()
	out := filepath.Join(dir, "result.json")

	opts := DefaultPatchOptions()
	opts.Enabled = true
	opts.OutputPath = out

	_, err := GeneratePatch(changes, opts)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	var reloaded Patch
	if err := json.Unmarshal(data, &reloaded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if len(reloaded.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(reloaded.Entries))
	}
	if reloaded.GeneratedAt == "" {
		t.Error("GeneratedAt should not be empty")
	}
}
