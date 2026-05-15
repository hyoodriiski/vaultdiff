package diff

import (
	"testing"
)

func TestPipelineWithHighlight(t *testing.T) {
	secretA := map[string]interface{}{
		"api_url":  "https://api.staging.example.com",
		"timeout":  "30s",
		"password": "hunter2",
	}
	secretB := map[string]interface{}{
		"api_url":  "https://api.prod.example.com",
		"timeout":  "30s",
		"password": "correct-horse",
	}

	changes := Compare(secretA, secretB)

	highlightOpts := HighlightOptions{Enabled: true, Prefix: "[[", Suffix: "]]"}
	changes = ApplyHighlight(changes, highlightOpts)

	filterOpts := DefaultFilterOptions()
	filterOpts.Types = []ChangeType{ChangeTypeModified}
	changes = Filter(changes, filterOpts)

	if len(changes) != 2 {
		t.Fatalf("expected 2 modified changes, got %d", len(changes))
	}

	for _, c := range changes {
		if c.Type != ChangeTypeModified {
			t.Errorf("expected only Modified changes, got %v for key %q", c.Type, c.Key)
		}
	}
}

func TestPipelineHighlight_UnchangedNotMarked(t *testing.T) {
	secretA := map[string]interface{}{
		"stable": "value",
		"changed": "before",
	}
	secretB := map[string]interface{}{
		"stable": "value",
		"changed": "after",
	}

	changes := Compare(secretA, secretB)
	opts := HighlightOptions{Enabled: true, Prefix: "[[", Suffix: "]]"}
	result := ApplyHighlight(changes, opts)

	for _, c := range result {
		if c.Key == "stable" && c.Type == ChangeTypeUnchanged {
			if c.OldValue != "value" || c.NewValue != "value" {
				t.Errorf("unchanged key should not be highlighted")
			}
		}
	}
}
