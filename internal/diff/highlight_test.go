package diff

import (
	"testing"
)

func makeHighlightChanges() []Change {
	return []Change{
		{Key: "host", Type: ChangeTypeModified, OldValue: "localhost", NewValue: "remotehost"},
		{Key: "port", Type: ChangeTypeModified, OldValue: "8080", NewValue: "9090"},
		{Key: "token", Type: ChangeTypeAdded, NewValue: "abc123"},
		{Key: "debug", Type: ChangeTypeUnchanged, OldValue: "true", NewValue: "true"},
	}
}

func TestApplyHighlight_Disabled(t *testing.T) {
	changes := makeHighlightChanges()
	opts := DefaultHighlightOptions()
	opts.Enabled = false

	result := ApplyHighlight(changes, opts)
	if result[0].OldValue != "localhost" {
		t.Errorf("expected unchanged OldValue, got %q", result[0].OldValue)
	}
}

func TestApplyHighlight_OnlyModified(t *testing.T) {
	changes := makeHighlightChanges()
	opts := DefaultHighlightOptions()
	opts.Enabled = true

	result := ApplyHighlight(changes, opts)

	// Added and Unchanged should be untouched.
	if result[2].NewValue != "abc123" {
		t.Errorf("added key should be unchanged, got %q", result[2].NewValue)
	}
	if result[3].OldValue != "true" {
		t.Errorf("unchanged key should be untouched, got %q", result[3].OldValue)
	}
}

func TestApplyHighlight_MarksChangedChars(t *testing.T) {
	changes := []Change{
		{Key: "env", Type: ChangeTypeModified, OldValue: "staging", NewValue: "production"},
	}
	opts := HighlightOptions{Enabled: true, Prefix: "[[", Suffix: "]]"}

	result := ApplyHighlight(changes, opts)

	if result[0].OldValue == "staging" {
		t.Error("expected OldValue to be highlighted")
	}
	if result[0].NewValue == "production" {
		t.Error("expected NewValue to be highlighted")
	}
}

func TestApplyHighlight_IdenticalValues(t *testing.T) {
	changes := []Change{
		{Key: "x", Type: ChangeTypeModified, OldValue: "same", NewValue: "same"},
	}
	opts := HighlightOptions{Enabled: true, Prefix: "[[", Suffix: "]]"}
	result := ApplyHighlight(changes, opts)

	if result[0].OldValue != "same" {
		t.Errorf("identical values should not be highlighted, got %q", result[0].OldValue)
	}
}

func TestParseHighlightFlags_Default(t *testing.T) {
	opts, err := ParseHighlightFlags(true, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Prefix != "[[" || opts.Suffix != "]]" {
		t.Errorf("expected default markers, got %q %q", opts.Prefix, opts.Suffix)
	}
}

func TestParseHighlightFlags_CustomMarkers(t *testing.T) {
	opts, err := ParseHighlightFlags(true, ">>:<<")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Prefix != ">>" || opts.Suffix != "<<" {
		t.Errorf("expected custom markers, got %q %q", opts.Prefix, opts.Suffix)
	}
}

func TestParseHighlightFlags_InvalidMarkers(t *testing.T) {
	_, err := ParseHighlightFlags(true, "nocolon")
	if err == nil {
		t.Error("expected error for missing colon separator")
	}
}

func TestParseHighlightFlags_EmptyPart(t *testing.T) {
	_, err := ParseHighlightFlags(true, ":suffix")
	if err == nil {
		t.Error("expected error for empty prefix")
	}
}
