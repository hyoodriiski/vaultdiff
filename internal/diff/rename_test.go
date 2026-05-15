package diff

import (
	"testing"
)

func makeRenameChanges() []Change {
	return []Change{
		{Key: "old_password", Type: ChangeTypeRemoved, OldValue: "secret123"},
		{Key: "new_password", Type: ChangeTypeAdded, NewValue: "secret123"},
		{Key: "host", Type: ChangeTypeUnchanged, OldValue: "localhost", NewValue: "localhost"},
		{Key: "orphan_removed", Type: ChangeTypeRemoved, OldValue: "totally-different"},
	}
}

func TestApplyRename_Disabled(t *testing.T) {
	changes := makeRenameChanges()
	opts := DefaultRenameOptions()
	opts.Enabled = false
	result := ApplyRename(changes, opts)
	for _, c := range result {
		if c.Meta != nil && (c.Meta["renamed-to"] != "" || c.Meta["renamed-from"] != "") {
			t.Errorf("expected no rename metadata when disabled, got %v", c.Meta)
		}
	}
}

func TestApplyRename_DetectsExactMatch(t *testing.T) {
	changes := makeRenameChanges()
	opts := DefaultRenameOptions()
	opts.Enabled = true
	opts.Similarity = 1.0
	result := ApplyRename(changes, opts)

	var removedMeta, addedMeta map[string]string
	for _, c := range result {
		if c.Key == "old_password" {
			removeMeta := c.Meta
			removeMeta = removeMeta
			removeMeta = c.Meta
			removedMeta = removeMeta
		}
		if c.Key == "new_password" {
			addedMeta = c.Meta
		}
	}
	if removedMeta == nil || removedMeta["renamed-to"] != "new_password" {
		t.Errorf("expected renamed-to=new_password on old_password, got %v", removedMeta)
	}
	if addedMeta == nil || addedMeta["renamed-from"] != "old_password" {
		t.Errorf("expected renamed-from=old_password on new_password, got %v", addedMeta)
	}
}

func TestApplyRename_OrphanNotAnnotated(t *testing.T) {
	changes := makeRenameChanges()
	opts := DefaultRenameOptions()
	opts.Enabled = true
	opts.Similarity = 0.8
	result := ApplyRename(changes, opts)
	for _, c := range result {
		if c.Key == "orphan_removed" {
			if c.Meta != nil && c.Meta["renamed-to"] != "" {
				t.Errorf("orphan should not have rename metadata, got %v", c.Meta)
			}
		}
	}
}

func TestValueSimilarity_Equal(t *testing.T) {
	if s := valueSimilarity("abc", "abc"); s != 1.0 {
		t.Errorf("expected 1.0, got %v", s)
	}
}

func TestValueSimilarity_Empty(t *testing.T) {
	if s := valueSimilarity("", "abc"); s != 0.0 {
		t.Errorf("expected 0.0, got %v", s)
	}
}

func TestParseRenameFlags_Valid(t *testing.T) {
	opts, err := ParseRenameFlags("true", "0.75")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected Enabled=true")
	}
	if opts.Similarity != 0.75 {
		t.Errorf("expected Similarity=0.75, got %v", opts.Similarity)
	}
}

func TestParseRenameFlags_InvalidSimilarity(t *testing.T) {
	_, err := ParseRenameFlags("true", "1.5")
	if err == nil {
		t.Error("expected error for similarity > 1.0")
	}
}

func TestParseRenameFlags_InvalidEnabled(t *testing.T) {
	_, err := ParseRenameFlags("maybe", "")
	if err == nil {
		t.Error("expected error for invalid enabled value")
	}
}
