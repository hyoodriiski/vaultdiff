package diff

import (
	"sort"
	"time"
)

// TimelineEntry records a snapshot of changes at a point in time.
type TimelineEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Changes   []Change  `json:"changes"`
	Label     string    `json:"label,omitempty"`
}

// TimelineOptions controls timeline recording behaviour.
type TimelineOptions struct {
	Enabled   bool
	MaxEvents int
	Label     string
}

// DefaultTimelineOptions returns sensible defaults.
func DefaultTimelineOptions() TimelineOptions {
	return TimelineOptions{
		Enabled:   false,
		MaxEvents: 100,
	}
}

// Timeline holds an ordered sequence of diff snapshots.
type Timeline struct {
	Entries []TimelineEntry
	opts    TimelineOptions
}

// NewTimeline creates an empty Timeline with the given options.
func NewTimeline(opts TimelineOptions) *Timeline {
	return &Timeline{opts: opts}
}

// Record appends a new entry to the timeline. If MaxEvents is exceeded the
// oldest entry is evicted.
func (t *Timeline) Record(changes []Change) {
	if !t.opts.Enabled {
		return
	}
	entry := TimelineEntry{
		Timestamp: time.Now().UTC(),
		Changes:   changes,
		Label:     t.opts.Label,
	}
	t.Entries = append(t.Entries, entry)
	if t.opts.MaxEvents > 0 && len(t.Entries) > t.opts.MaxEvents {
		t.Entries = t.Entries[len(t.Entries)-t.opts.MaxEvents:]
	}
}

// Since returns all entries recorded at or after the given time.
func (t *Timeline) Since(since time.Time) []TimelineEntry {
	var out []TimelineEntry
	for _, e := range t.Entries {
		if !e.Timestamp.Before(since) {
			out = append(out, e)
		}
	}
	return out
}

// SortedByTime returns entries in ascending timestamp order.
func (t *Timeline) SortedByTime() []TimelineEntry {
	copy := append([]TimelineEntry(nil), t.Entries...)
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].Timestamp.Before(copy[j].Timestamp)
	})
	return copy
}

// Len returns the number of recorded entries.
func (t *Timeline) Len() int { return len(t.Entries) }
