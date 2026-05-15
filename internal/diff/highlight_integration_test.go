package diff

import (
	"strings"
	"testing"
)

func TestHighlightAfterCompare(t *testing.T) {
	secretA := map[string]interface{}{
		"db_host": "db.staging.internal",
		"db_port": "5432",
	}
	secretB := map[string]interface{}{
		"db_host": "db.production.internal",
		"db_port": "5432",
	}

	changes := Compare(secretA, secretB)
	opts := HighlightOptions{Enabled: true, Prefix: "[[", Suffix: "]]"}
	highlighted := ApplyHighlight(changes, opts)

	var modifiedChange *Change
	for i := range highlighted {
		if highlighted[i].Key == "db_host" {
			modifiedChange = &highlighted[i]
			break
		}
	}

	if modifiedChange == nil {
		t.Fatal("expected db_host to be in changes")
	}
	if modifiedChange.Type != ChangeTypeModified {
		t.Fatalf("expected Modified, got %v", modifiedChange.Type)
	}
	if !strings.Contains(modifiedChange.OldValue, "[[") {
		t.Errorf("expected OldValue to contain highlight markers, got %q", modifiedChange.OldValue)
	}
	if !strings.Contains(modifiedChange.NewValue, "[[") {
		t.Errorf("expected NewValue to contain highlight markers, got %q", modifiedChange.NewValue)
	}
}

func TestHighlightDisabledPreservesValues(t *testing.T) {
	secretA := map[string]interface{}{"key": "old_value"}
	secretB := map[string]interface{}{"key": "new_value"}

	changes := Compare(secretA, secretB)
	opts := DefaultHighlightOptions() // Enabled = false
	highlighted := ApplyHighlight(changes, opts)

	for _, c := range highlighted {
		if c.Key == "key" {
			if c.OldValue != "old_value" {
				t.Errorf("expected original OldValue, got %q", c.OldValue)
			}
			if c.NewValue != "new_value" {
				t.Errorf("expected original NewValue, got %q", c.NewValue)
			}
			return
		}
	}
	t.Error("key not found in changes")
}
