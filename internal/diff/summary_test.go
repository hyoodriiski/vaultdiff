package diff_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func makeSummaryChanges() []diff.Change {
	return []diff.Change{
		{Key: "a", Type: diff.Added},
		{Key: "b", Type: diff.Added},
		{Key: "c", Type: diff.Removed},
		{Key: "d", Type: diff.Modified},
		{Key: "e", Type: diff.Unchanged},
	}
}

func TestSummarize_Counts(t *testing.T) {
	changes := makeSummaryChanges()
	s := diff.Summarize(changes)

	if s.Total != 5 {
		t.Errorf("Total: got %d, want 5", s.Total)
	}
	if s.Added != 2 {
		t.Errorf("Added: got %d, want 2", s.Added)
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

func TestSummarize_Empty(t *testing.T) {
	s := diff.Summarize([]diff.Change{})
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
	if s.Percent(0) != 0 {
		t.Error("Percent on empty summary should be 0")
	}
}

func TestSummary_Percent(t *testing.T) {
	s := diff.Summary{Total: 4, Added: 2}
	got := s.Percent(s.Added)
	if got != 50.0 {
		t.Errorf("Percent: got %.1f, want 50.0", got)
	}
}

func TestParseSummaryFlags_Valid(t *testing.T) {
	opts, err := diff.ParseSummaryFlags("counts,percents")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.ShowCounts || !opts.ShowPercents {
		t.Error("expected both ShowCounts and ShowPercents to be true")
	}
}

func TestParseSummaryFlags_Invalid(t *testing.T) {
	_, err := diff.ParseSummaryFlags("unknown")
	if err == nil {
		t.Error("expected error for unknown flag, got nil")
	}
}

func TestParseSummaryFlags_Empty(t *testing.T) {
	opts, err := diff.ParseSummaryFlags("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.ShowCounts {
		t.Error("default options should have ShowCounts=true")
	}
}
