package diff

import (
	"testing"
)

func makeGroupChanges() []Change {
	return []Change{
		{Key: "db:host", Type: ChangeTypeAdded, NewValue: "localhost"},
		{Key: "db:port", Type: ChangeTypeModified, OldValue: "5432", NewValue: "5433"},
		{Key: "app:name", Type: ChangeTypeUnchanged, OldValue: "myapp", NewValue: "myapp"},
		{Key: "app:version", Type: ChangeTypeRemoved, OldValue: "1.0"},
		{Key: "token", Type: ChangeTypeAdded, NewValue: "abc123"},
	}
}

func TestGroup_Disabled(t *testing.T) {
	changes := makeGroupChanges()
	opts := DefaultGroupOptions()
	groups := Group(changes, opts)

	if len(groups) != 1 {
		t.Fatalf("expected 1 group when disabled, got %d", len(groups))
	}
	if groups[0].Label != "" {
		t.Errorf("expected empty label, got %q", groups[0].Label)
	}
	if len(groups[0].Changes) != len(changes) {
		t.Errorf("expected %d changes, got %d", len(changes), len(groups[0].Changes))
	}
}

func TestGroup_ByType(t *testing.T) {
	changes := makeGroupChanges()
	opts := GroupOptions{Enabled: true, GroupBy: "type", Separator: ":"}
	groups := Group(changes, opts)

	if len(groups) == 0 {
		t.Fatal("expected at least one group")
	}

	total := 0
	for _, g := range groups {
		total += len(g.Changes)
		for _, c := range g.Changes {
			if string(c.Type) != g.Label {
				t.Errorf("change type %q does not match group label %q", c.Type, g.Label)
			}
		}
	}
	if total != len(changes) {
		t.Errorf("expected %d total changes across groups, got %d", len(changes), total)
	}
}

func TestGroup_ByPrefix(t *testing.T) {
	changes := makeGroupChanges()
	opts := GroupOptions{Enabled: true, GroupBy: "prefix", Separator: ":"}
	groups := Group(changes, opts)

	labels := map[string]bool{}
	for _, g := range groups {
		labels[g.Label] = true
	}

	if !labels["db"] {
		t.Error("expected group with label 'db'")
	}
	if !labels["app"] {
		t.Error("expected group with label 'app'")
	}
	if !labels["token"] {
		t.Error("expected group with label 'token' (no separator)")
	}
}

func TestGroup_Empty(t *testing.T) {
	opts := GroupOptions{Enabled: true, GroupBy: "type", Separator: ":"}
	groups := Group([]Change{}, opts)

	if len(groups) != 1 {
		t.Fatalf("expected 1 group for empty input, got %d", len(groups))
	}
	if len(groups[0].Changes) != 0 {
		t.Errorf("expected 0 changes in group, got %d", len(groups[0].Changes))
	}
}

func TestParseGroupFlags_Valid(t *testing.T) {
	opts, err := ParseGroupFlags(true, "prefix", "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.GroupBy != "prefix" {
		t.Errorf("expected GroupBy=prefix, got %q", opts.GroupBy)
	}
	if opts.Separator != "/" {
		t.Errorf("expected Separator='/', got %q", opts.Separator)
	}
}

func TestParseGroupFlags_InvalidGroupBy(t *testing.T) {
	_, err := ParseGroupFlags(true, "invalid", ":")
	if err == nil {
		t.Error("expected error for invalid group-by value")
	}
}

func TestParseGroupFlags_Disabled(t *testing.T) {
	opts, err := ParseGroupFlags(false, "invalid", "")
	if err != nil {
		t.Fatalf("expected no error when disabled, got %v", err)
	}
	if opts.Enabled {
		t.Error("expected Enabled=false")
	}
}
