package diff

import (
	"strings"
	"testing"
)

func makeAnnotateChanges() []Change {
	return []Change{
		{Key: "alpha", Type: ChangeAdded, NewValue: "new"},
		{Key: "beta", Type: ChangeRemoved, OldValue: "old"},
		{Key: "gamma", Type: ChangeModified, OldValue: "v1", NewValue: "v2"},
		{Key: "delta", Type: ChangeUnchanged, OldValue: "same", NewValue: "same"},
	}
}

func TestApplyAnnotate_Disabled(t *testing.T) {
	changes := makeAnnotateChanges()
	opts := DefaultAnnotateOptions()
	result := ApplyAnnotate(changes, opts)
	if result[0].Key != "alpha" {
		t.Errorf("expected key 'alpha', got %q", result[0].Key)
	}
}

func TestApplyAnnotate_ShowIndex(t *testing.T) {
	changes := makeAnnotateChanges()
	opts := AnnotateOptions{Enabled: true, ShowIndex: true}
	result := ApplyAnnotate(changes, opts)
	if !strings.Contains(result[0].Key, "[1]") {
		t.Errorf("expected index annotation in key, got %q", result[0].Key)
	}
	if !strings.Contains(result[3].Key, "[4]") {
		t.Errorf("expected index 4 in last key, got %q", result[3].Key)
	}
}

func TestApplyAnnotate_ShowSource_Added(t *testing.T) {
	changes := makeAnnotateChanges()
	opts := AnnotateOptions{Enabled: true, ShowSource: true, SourceA: "left", SourceB: "right"}
	result := ApplyAnnotate(changes, opts)
	if !strings.Contains(result[0].NewValue, "from right") {
		t.Errorf("expected source annotation in NewValue, got %q", result[0].NewValue)
	}
}

func TestApplyAnnotate_ShowSource_Removed(t *testing.T) {
	changes := makeAnnotateChanges()
	opts := AnnotateOptions{Enabled: true, ShowSource: true, SourceA: "left", SourceB: "right"}
	result := ApplyAnnotate(changes, opts)
	if !strings.Contains(result[1].OldValue, "from left") {
		t.Errorf("expected source annotation in OldValue, got %q", result[1].OldValue)
	}
}

func TestApplyAnnotate_ShowSource_Modified(t *testing.T) {
	changes := makeAnnotateChanges()
	opts := AnnotateOptions{Enabled: true, ShowSource: true, SourceA: "src-a", SourceB: "src-b"}
	result := ApplyAnnotate(changes, opts)
	if !strings.Contains(result[2].OldValue, "from src-a") {
		t.Errorf("expected OldValue annotation, got %q", result[2].OldValue)
	}
	if !strings.Contains(result[2].NewValue, "from src-b") {
		t.Errorf("expected NewValue annotation, got %q", result[2].NewValue)
	}
}

func TestApplyAnnotate_UnchangedNotAnnotatedSource(t *testing.T) {
	changes := makeAnnotateChanges()
	opts := AnnotateOptions{Enabled: true, ShowSource: true, SourceA: "a", SourceB: "b"}
	result := ApplyAnnotate(changes, opts)
	// Unchanged entries should not have source annotations on values
	if strings.Contains(result[3].OldValue, "from") {
		t.Errorf("unchanged entry should not be source-annotated, got %q", result[3].OldValue)
	}
}

func TestParseAnnotateFlags_Defaults(t *testing.T) {
	opts, err := ParseAnnotateFlags(false, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.SourceA != "a" || opts.SourceB != "b" {
		t.Errorf("unexpected default sources: %q %q", opts.SourceA, opts.SourceB)
	}
}

func TestParseAnnotateFlags_SourceOverride(t *testing.T) {
	opts, err := ParseAnnotateFlags(false, true, "a=vault/prod,b=vault/staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.SourceA != "vault/prod" {
		t.Errorf("expected SourceA=vault/prod, got %q", opts.SourceA)
	}
	if opts.SourceB != "vault/staging" {
		t.Errorf("expected SourceB=vault/staging, got %q", opts.SourceB)
	}
}
