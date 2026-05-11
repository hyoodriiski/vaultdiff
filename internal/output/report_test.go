package output_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/output"
)

func makeChanges() []diff.Change {
	return []diff.Change{
		{Key: "db_host", Type: diff.Added, NewValue: "localhost"},
		{Key: "db_pass", Type: diff.Removed, OldValue: "secret"},
		{Key: "db_user", Type: diff.Modified, OldValue: "root", NewValue: "admin"},
		{Key: "app_env", Type: diff.Unchanged, OldValue: "prod", NewValue: "prod"},
	}
}

func TestNewReport_Summary(t *testing.T) {
	changes := makeChanges()
	r := output.NewReport("secret/src", "secret/dst", changes)

	if r.Summary.Added != 1 {
		t.Errorf("expected Added=1, got %d", r.Summary.Added)
	}
	if r.Summary.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", r.Summary.Removed)
	}
	if r.Summary.Modified != 1 {
		t.Errorf("expected Modified=1, got %d", r.Summary.Modified)
	}
	if r.Summary.Unchanged != 1 {
		t.Errorf("expected Unchanged=1, got %d", r.Summary.Unchanged)
	}
	if r.Summary.Total != 4 {
		t.Errorf("expected Total=4, got %d", r.Summary.Total)
	}
}

func TestNewReport_Paths(t *testing.T) {
	r := output.NewReport("secret/a", "secret/b", nil)
	if r.SourcePath != "secret/a" || r.TargetPath != "secret/b" {
		t.Errorf("unexpected paths: %q %q", r.SourcePath, r.TargetPath)
	}
}

func TestWriteJSON_ValidOutput(t *testing.T) {
	changes := makeChanges()
	r := output.NewReport("secret/src", "secret/dst", changes)

	var buf bytes.Buffer
	if err := output.WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}

	var decoded output.Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to decode JSON output: %v", err)
	}
	if decoded.Summary.Total != 4 {
		t.Errorf("decoded summary total mismatch: got %d", decoded.Summary.Total)
	}
}

func TestNewReport_EmptyChanges(t *testing.T) {
	r := output.NewReport("a", "b", []diff.Change{})
	if r.Summary.Total != 0 {
		t.Errorf("expected 0 total, got %d", r.Summary.Total)
	}
}
