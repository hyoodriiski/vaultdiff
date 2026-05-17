package diff

import (
	"strings"
	"testing"
)

func makeEnrichChanges() []Change {
	return []Change{
		{Key: "db/password", Type: ChangeTypeAdded, NewValue: "s3cr3t"},
		{Key: "api_key", Type: ChangeTypeRemoved, OldValue: "old-key"},
		{Key: "host", Type: ChangeTypeModified, OldValue: "localhost", NewValue: "prod.example.com"},
		{Key: "port", Type: ChangeTypeUnchanged, OldValue: "5432", NewValue: "5432"},
	}
}

func TestApplyEnrich_Disabled(t *testing.T) {
	changes := makeEnrichChanges()
	opts := DefaultEnrichOptions()
	opts.Enabled = false
	out := ApplyEnrich(changes, opts)
	for _, c := range out {
		if len(c.Annotations) != 0 {
			t.Errorf("expected no annotations when disabled, got %v", c.Annotations)
		}
	}
}

func TestApplyEnrich_KeyLength(t *testing.T) {
	changes := makeEnrichChanges()
	opts := DefaultEnrichOptions()
	opts.Enabled = true
	opts.AddValueLength = false
	opts.AddChangeID = false
	out := ApplyEnrich(changes, opts)
	for _, c := range out {
		found := false
		for _, a := range c.Annotations {
			if strings.HasPrefix(a, "key_len=") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected key_len annotation for key %q", c.Key)
		}
	}
}

func TestApplyEnrich_ChangeID(t *testing.T) {
	changes := []Change{
		{Key: "my/key", Type: ChangeTypeAdded, NewValue: "val"},
	}
	opts := DefaultEnrichOptions()
	opts.Enabled = true
	opts.AddKeyLength = false
	opts.AddValueLength = false
	opts.IDPrefix = "vd"
	out := ApplyEnrich(changes, opts)
	found := false
	for _, a := range out[0].Annotations {
		if a == "id=vd_my_key_added" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected change ID annotation, got %v", out[0].Annotations)
	}
}

func TestApplyEnrich_AddedHasNoOldLen(t *testing.T) {
	changes := []Change{
		{Key: "token", Type: ChangeTypeAdded, NewValue: "abc123"},
	}
	opts := DefaultEnrichOptions()
	opts.Enabled = true
	opts.AddChangeID = false
	opts.AddKeyLength = false
	out := ApplyEnrich(changes, opts)
	for _, a := range out[0].Annotations {
		if strings.HasPrefix(a, "old_len=") {
			t.Errorf("added change should not have old_len annotation")
		}
	}
}

func TestApplyEnrich_NoDuplicateAnnotations(t *testing.T) {
	changes := []Change{
		{Key: "k", Type: ChangeTypeModified, OldValue: "a", NewValue: "b", Annotations: []string{"key_len=1"}},
	}
	opts := DefaultEnrichOptions()
	opts.Enabled = true
	opts.AddValueLength = false
	opts.AddChangeID = false
	out := ApplyEnrich(changes, opts)
	count := 0
	for _, a := range out[0].Annotations {
		if strings.HasPrefix(a, "key_len=") {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 key_len annotation, got %d", count)
	}
}
