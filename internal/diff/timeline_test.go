package diff

import (
	"testing"
	"time"
)

func makeTimelineChanges() []Change {
	return []Change{
		{Key: "db_pass", OldValue: "old", NewValue: "new", Type: ChangeTypeModified},
		{Key: "api_key", OldValue: "", NewValue: "abc", Type: ChangeTypeAdded},
	}
}

func TestNewTimeline_Disabled(t *testing.T) {
	opts := DefaultTimelineOptions()
	tl := NewTimeline(opts)
	tl.Record(makeTimelineChanges())
	if tl.Len() != 0 {
		t.Fatalf("expected 0 entries when disabled, got %d", tl.Len())
	}
}

func TestTimeline_RecordsEntry(t *testing.T) {
	opts := DefaultTimelineOptions()
	opts.Enabled = true
	tl := NewTimeline(opts)
	tl.Record(makeTimelineChanges())
	if tl.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", tl.Len())
	}
	if len(tl.Entries[0].Changes) != 2 {
		t.Errorf("expected 2 changes in entry, got %d", len(tl.Entries[0].Changes))
	}
}

func TestTimeline_MaxEventsEvictsOldest(t *testing.T) {
	opts := DefaultTimelineOptions()
	opts.Enabled = true
	opts.MaxEvents = 3
	tl := NewTimeline(opts)
	for i := 0; i < 5; i++ {
		tl.Record(makeTimelineChanges())
	}
	if tl.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", tl.Len())
	}
}

func TestTimeline_Since(t *testing.T) {
	opts := DefaultTimelineOptions()
	opts.Enabled = true
	tl := NewTimeline(opts)

	past := time.Now().UTC().Add(-2 * time.Hour)
	tl.Entries = []TimelineEntry{
		{Timestamp: past, Changes: makeTimelineChanges()},
		{Timestamp: time.Now().UTC(), Changes: makeTimelineChanges()},
	}

	cutoff := time.Now().UTC().Add(-1 * time.Hour)
	result := tl.Since(cutoff)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry since cutoff, got %d", len(result))
	}
}

func TestTimeline_SortedByTime(t *testing.T) {
	opts := DefaultTimelineOptions()
	opts.Enabled = true
	tl := NewTimeline(opts)

	now := time.Now().UTC()
	tl.Entries = []TimelineEntry{
		{Timestamp: now.Add(2 * time.Minute), Label: "third"},
		{Timestamp: now, Label: "first"},
		{Timestamp: now.Add(1 * time.Minute), Label: "second"},
	}

	sorted := tl.SortedByTime()
	if sorted[0].Label != "first" || sorted[1].Label != "second" || sorted[2].Label != "third" {
		t.Errorf("unexpected sort order: %v %v %v", sorted[0].Label, sorted[1].Label, sorted[2].Label)
	}
}

func TestTimeline_LabelPropagated(t *testing.T) {
	opts := DefaultTimelineOptions()
	opts.Enabled = true
	opts.Label = "prod-snapshot"
	tl := NewTimeline(opts)
	tl.Record(makeTimelineChanges())
	if tl.Entries[0].Label != "prod-snapshot" {
		t.Errorf("expected label 'prod-snapshot', got %q", tl.Entries[0].Label)
	}
}
