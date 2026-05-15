package diff

import "testing"

func TestRenameAfterCompare(t *testing.T) {
	secretA := map[string]interface{}{
		"api_key":  "hunter2",
		"endpoint": "https://example.com",
	}
	secretB := map[string]interface{}{
		"token":    "hunter2",
		"endpoint": "https://example.com",
	}

	changes := Compare(secretA, secretB)
	opts := DefaultRenameOptions()
	opts.Enabled = true
	opts.Similarity = 1.0
	result := ApplyRename(changes, opts)

	var apiKeyMeta, tokenMeta map[string]string
	for _, c := range result {
		switch c.Key {
		case "api_key":
			apiKeyMeta = c.Meta
		case "token":
			tokenMeta = c.Meta
		}
	}

	if apiKeyMeta == nil || apiKeyMeta["renamed-to"] != "token" {
		t.Errorf("expected api_key renamed-to=token, got %v", apiKeyMeta)
	}
	if tokenMeta == nil || tokenMeta["renamed-from"] != "api_key" {
		t.Errorf("expected token renamed-from=api_key, got %v", tokenMeta)
	}
}

func TestRenameDisabledPreservesChanges(t *testing.T) {
	secretA := map[string]interface{}{"x": "val"}
	secretB := map[string]interface{}{"y": "val"}

	changes := Compare(secretA, secretB)
	opts := DefaultRenameOptions()
	opts.Enabled = false
	result := ApplyRename(changes, opts)

	if len(result) != len(changes) {
		t.Errorf("expected same number of changes, got %d vs %d", len(result), len(changes))
	}
	for _, c := range result {
		if c.Meta != nil && (c.Meta["renamed-to"] != "" || c.Meta["renamed-from"] != "") {
			t.Errorf("unexpected rename metadata when disabled: %v", c.Meta)
		}
	}
}
