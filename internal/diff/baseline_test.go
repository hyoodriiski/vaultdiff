package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func makeBaselineChanges() []Change {
	return []Change{
		{Key: "host", Type: ChangeTypeUnchanged, OldValue: "localhost", NewValue: "localhost"},
		{Key: "password", Type: ChangeTypeModified, OldValue: "old-pass", NewValue: "new-pass"},
		{Key: "token", Type: ChangeTypeAdded, OldValue: "", NewValue: "abc123"},
		{Key: "legacy", Type: ChangeTypeRemoved, OldValue: "gone", NewValue: ""},
	}
}

func TestSaveBaseline_Disabled(t *testing.T) {
	opts := DefaultBaselineOptions()
	if err := SaveBaseline(makeBaselineChanges(), opts); err != nil {
		t.Fatalf("expected no error when disabled, got: %v", err)
	}
}

func TestSaveBaseline_WritesFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "baseline.json")
	opts := BaselineOptions{Enabled: true, FilePath: tmp}

	if err := SaveBaseline(makeBaselineChanges(), opts); err != nil {
		t.Fatalf("SaveBaseline error: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	var snap struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if snap.Data["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %q", snap.Data["host"])
	}
	if snap.Data["password"] != "new-pass" {
		t.Errorf("expected password=new-pass (new value), got %q", snap.Data["password"])
	}
	if snap.Data["token"] != "abc123" {
		t.Errorf("expected token=abc123, got %q", snap.Data["token"])
	}
	if snap.Data["legacy"] != "gone" {
		t.Errorf("expected legacy=gone (old value), got %q", snap.Data["legacy"])
	}
}

func TestLoadBaseline_ReturnsMap(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "baseline.json")
	opts := BaselineOptions{Enabled: true, FilePath: tmp}

	if err := SaveBaseline(makeBaselineChanges(), opts); err != nil {
		t.Fatalf("SaveBaseline error: %v", err)
	}

	m, err := LoadBaseline(tmp)
	if err != nil {
		t.Fatalf("LoadBaseline error: %v", err)
	}

	if m["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", m["host"])
	}
	if m["token"] != "abc123" {
		t.Errorf("expected token=abc123, got %v", m["token"])
	}
}

func TestLoadBaseline_EmptyPath(t *testing.T) {
	_, err := LoadBaseline("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
