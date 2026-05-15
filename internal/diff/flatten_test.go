package diff

import (
	"testing"
)

func makeFlattenChanges() []Change {
	return []Change{
		{Key: "config", OldValue: "map[host:localhost port:5432]", NewValue: "map[host:db.prod port:5432]", Type: ChangeTypeModified},
		{Key: "simple", OldValue: "alpha", NewValue: "beta", Type: ChangeTypeModified},
		{Key: "added_map", OldValue: "", NewValue: "map[token:abc]", Type: ChangeTypeAdded},
	}
}

func TestApplyFlatten_Disabled(t *testing.T) {
	changes := makeFlattenChanges()
	opts := DefaultFlattenOptions()
	opts.Enabled = false

	result := ApplyFlatten(changes, opts)
	if len(result) != len(changes) {
		t.Fatalf("expected %d changes, got %d", len(changes), len(result))
	}
	if result[0].Key != "config" {
		t.Errorf("expected key 'config', got %q", result[0].Key)
	}
}

func TestApplyFlatten_ExpandsNestedMap(t *testing.T) {
	changes := []Change{
		{Key: "db", OldValue: "map[host:localhost port:5432]", NewValue: "map[host:db.prod port:5432]", Type: ChangeTypeModified},
	}
	opts := DefaultFlattenOptions()
	opts.Enabled = true

	result := ApplyFlatten(changes, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 flattened changes, got %d", len(result))
	}

	keys := map[string]bool{}
	for _, c := range result {
		keys[c.Key] = true
	}
	if !keys["db.host"] {
		t.Error("expected key 'db.host'")
	}
	if !keys["db.port"] {
		t.Error("expected key 'db.port'")
	}
}

func TestApplyFlatten_PassesThroughSimple(t *testing.T) {
	changes := []Change{
		{Key: "simple", OldValue: "alpha", NewValue: "beta", Type: ChangeTypeModified},
	}
	opts := DefaultFlattenOptions()
	opts.Enabled = true

	result := ApplyFlatten(changes, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if result[0].Key != "simple" {
		t.Errorf("expected key 'simple', got %q", result[0].Key)
	}
}

func TestApplyFlatten_AddedMapInfersType(t *testing.T) {
	changes := []Change{
		{Key: "creds", OldValue: "", NewValue: "map[token:abc]", Type: ChangeTypeAdded},
	}
	opts := DefaultFlattenOptions()
	opts.Enabled = true

	result := ApplyFlatten(changes, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if result[0].Type != ChangeTypeAdded {
		t.Errorf("expected ChangeTypeAdded, got %v", result[0].Type)
	}
	if result[0].Key != "creds.token" {
		t.Errorf("expected key 'creds.token', got %q", result[0].Key)
	}
}

func TestApplyFlatten_CustomSeparator(t *testing.T) {
	changes := []Change{
		{Key: "app", OldValue: "map[env:dev]", NewValue: "map[env:prod]", Type: ChangeTypeModified},
	}
	opts := DefaultFlattenOptions()
	opts.Enabled = true
	opts.Separator = "/"

	result := ApplyFlatten(changes, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if result[0].Key != "app/env" {
		t.Errorf("expected key 'app/env', got %q", result[0].Key)
	}
}
