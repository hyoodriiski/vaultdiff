package diff_test

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
)

func TestCompare_Added(t *testing.T) {
	a := map[string]interface{}{}
	b := map[string]interface{}{"newkey": "newval"}

	result := diff.Compare("secret/a", "secret/b", a, b)

	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Added {
		t.Errorf("expected Added, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].Key != "newkey" {
		t.Errorf("expected key 'newkey', got %s", result.Changes[0].Key)
	}
}

func TestCompare_Removed(t *testing.T) {
	a := map[string]interface{}{"oldkey": "oldval"}
	b := map[string]interface{}{}

	result := diff.Compare("secret/a", "secret/b", a, b)

	if result.Changes[0].Type != diff.Removed {
		t.Errorf("expected Removed, got %s", result.Changes[0].Type)
	}
}

func TestCompare_Modified(t *testing.T) {
	a := map[string]interface{}{"key": "val1"}
	b := map[string]interface{}{"key": "val2"}

	result := diff.Compare("secret/a", "secret/b", a, b)

	if result.Changes[0].Type != diff.Modified {
		t.Errorf("expected Modified, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].OldValue != "val1" {
		t.Errorf("expected OldValue 'val1', got %v", result.Changes[0].OldValue)
	}
	if result.Changes[0].NewValue != "val2" {
		t.Errorf("expected NewValue 'val2', got %v", result.Changes[0].NewValue)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	a := map[string]interface{}{"key": "same"}
	b := map[string]interface{}{"key": "same"}

	result := diff.Compare("secret/a", "secret/b", a, b)

	if result.Changes[0].Type != diff.Unchanged {
		t.Errorf("expected Unchanged, got %s", result.Changes[0].Type)
	}
	if result.HasChanges() {
		t.Error("expected HasChanges() to be false")
	}
}

func TestCompare_Mixed(t *testing.T) {
	a := map[string]interface{}{"keep": "v", "remove": "x", "change": "old"}
	b := map[string]interface{}{"keep": "v", "add": "y", "change": "new"}

	result := diff.Compare("secret/a", "secret/b", a, b)

	if !result.HasChanges() {
		t.Error("expected HasChanges() to be true")
	}

	types := map[string]diff.ChangeType{}
	for _, c := range result.Changes {
		types[c.Key] = c.Type
	}

	if types["keep"] != diff.Unchanged {
		t.Errorf("keep: expected Unchanged, got %s", types["keep"])
	}
	if types["remove"] != diff.Removed {
		t.Errorf("remove: expected Removed, got %s", types["remove"])
	}
	if types["add"] != diff.Added {
		t.Errorf("add: expected Added, got %s", types["add"])
	}
	if types["change"] != diff.Modified {
		t.Errorf("change: expected Modified, got %s", types["change"])
	}
}
