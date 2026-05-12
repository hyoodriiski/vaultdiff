package diff

import (
	"strings"
	"testing"
)

func TestAnnotateAfterCompare(t *testing.T) {
	secretA := map[string]interface{}{
		"host":     "db.prod",
		"password": "secret",
	}
	secretB := map[string]interface{}{
		"host":    "db.staging",
		"api_key": "xyz",
	}

	changes := Compare(secretA, secretB)
	opts := AnnotateOptions{
		Enabled:    true,
		ShowSource: true,
		SourceA:    "prod",
		SourceB:    "staging",
	}
	annotated := ApplyAnnotate(changes, opts)

	for _, c := range annotated {
		switch c.Type {
		case ChangeAdded:
			if !strings.Contains(c.NewValue, "from staging") {
				t.Errorf("added key %q: expected 'from staging' in NewValue, got %q", c.Key, c.NewValue)
			}
		case ChangeRemoved:
			if !strings.Contains(c.OldValue, "from prod") {
				t.Errorf("removed key %q: expected 'from prod' in OldValue, got %q", c.Key, c.OldValue)
			}
		case ChangeModified:
			if !strings.Contains(c.OldValue, "from prod") {
				t.Errorf("modified key %q: expected 'from prod' in OldValue", c.Key)
			}
			if !strings.Contains(c.NewValue, "from staging") {
				t.Errorf("modified key %q: expected 'from staging' in NewValue", c.Key)
			}
		}
	}
}

func TestAnnotateWithIndexAfterCompare(t *testing.T) {
	secretA := map[string]interface{}{"x": "1", "y": "2"}
	secretB := map[string]interface{}{"x": "1", "z": "3"}

	changes := Compare(secretA, secretB)
	opts := AnnotateOptions{
		Enabled:   true,
		ShowIndex: true,
	}
	annotated := ApplyAnnotate(changes, opts)

	for i, c := range annotated {
		expected := strings.Contains(c.Key, "[")
		if !expected {
			t.Errorf("entry %d: expected index bracket in key %q", i, c.Key)
		}
	}
}
