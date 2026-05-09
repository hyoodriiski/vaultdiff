package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
)

func TestTextFormatter_NoChanges(t *testing.T) {
	result := diff.Compare("secret/a", "secret/b",
		map[string]interface{}{"key": "val"},
		map[string]interface{}{"key": "val"},
	)

	var buf bytes.Buffer
	if err := diff.TextFormatter(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No changes detected") {
		t.Errorf("expected 'No changes detected', got: %s", out)
	}
}

func TestTextFormatter_WithChanges(t *testing.T) {
	result := diff.Compare("secret/a", "secret/b",
		map[string]interface{}{"removed": "x", "changed": "old"},
		map[string]interface{}{"added": "y", "changed": "new"},
	)

	var buf bytes.Buffer
	if err := diff.TextFormatter(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "+ added") {
		t.Errorf("expected added key in output, got: %s", out)
	}
	if !strings.Contains(out, "- removed") {
		t.Errorf("expected removed key in output, got: %s", out)
	}
	if !strings.Contains(out, "~ changed") {
		t.Errorf("expected modified key in output, got: %s", out)
	}
}

func TestTextFormatter_Header(t *testing.T) {
	result := &diff.Result{
		PathA: "secret/env/prod",
		PathB: "secret/env/staging",
	}

	var buf bytes.Buffer
	diff.TextFormatter(&buf, result)

	out := buf.String()
	if !strings.Contains(out, "secret/env/prod") || !strings.Contains(out, "secret/env/staging") {
		t.Errorf("expected both paths in header, got: %s", out)
	}
}
