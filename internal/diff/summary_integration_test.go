package diff_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// TestSummarizeAfterCompare verifies Summarize works correctly on output
// produced by Compare, mirroring the pipeline used in production.
func TestSummarizeAfterCompare(t *testing.T) {
	old := map[string]interface{}{
		"host":     "old.example.com",
		"port":     "5432",
		"password": "secret",
	}
	new := map[string]interface{}{
		"host":    "new.example.com",
		"port":    "5432",
		"api_key": "abc123",
	}

	changes := diff.Compare(old, new)
	s := diff.Summarize(changes)

	if s.Total != 4 {
		t.Errorf("Total: got %d, want 4", s.Total)
	}
	if s.Added != 1 {
		t.Errorf("Added: got %d, want 1", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("Removed: got %d, want 1", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("Modified: got %d, want 1", s.Modified)
	}
	if s.Unchanged != 1 {
		t.Errorf("Unchanged: got %d, want 1", s.Unchanged)
	}
}

// TestSummarizeAfterFilter verifies totals reflect filtered-out entries.
func TestSummarizeAfterFilter(t *testing.T) {
	changes := []diff.Change{
		{Key: "a", Type: diff.Added},
		{Key: "b", Type: diff.Removed},
		{Key: "c", Type: diff.Unchanged},
	}

	fopts := diff.DefaultFilterOptions()
	fopts.IncludeUnchanged = false
	filtered := diff.Filter(changes, fopts)

	s := diff.Summarize(filtered)
	if s.Total != 2 {
		t.Errorf("Total after filter: got %d, want 2", s.Total)
	}
	if s.Unchanged != 0 {
		t.Errorf("Unchanged after filter: got %d, want 0", s.Unchanged)
	}
}
