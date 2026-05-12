package diff_test

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
)

// TestContextAfterCompare verifies that ApplyContext integrates correctly
// with the output of Compare when Enabled=true.
func TestContextAfterCompare(t *testing.T) {
	old := map[string]interface{}{
		"host":     "localhost",
		"port":     "5432",
		"user":     "admin",
		"password": "secret",
		"db":       "mydb",
	}
	new_ := map[string]interface{}{
		"host":     "localhost",
		"port":     "5432",
		"user":     "admin",
		"password": "n3wsecret", // modified
		"db":       "mydb",
	}

	changes := diff.Compare(old, new_)

	opts := diff.ContextOptions{Lines: 1, Enabled: true}
	result := diff.ApplyContext(changes, opts)

	// password is modified; with 1-line context we expect its immediate
	// sorted neighbours to appear too.  At minimum the modified key must
	// be present.
	found := false
	for _, c := range result {
		if c.Key == "password" && c.Type == diff.ChangeTypeModified {
			found = true
		}
	}
	if !found {
		t.Error("expected modified 'password' key in context output")
	}
	if len(result) == 0 {
		t.Error("result should not be empty")
	}
}

// TestContextDisabledPreservesAllChanges ensures that when context is
// disabled the full Compare result passes through unchanged.
func TestContextDisabledPreservesAllChanges(t *testing.T) {
	old := map[string]interface{}{"a": "1", "b": "2"}
	new_ := map[string]interface{}{"a": "1", "b": "3"}

	changes := diff.Compare(old, new_)
	opts := diff.ContextOptions{Enabled: false}
	result := diff.ApplyContext(changes, opts)

	if len(result) != len(changes) {
		t.Errorf("expected %d, got %d", len(changes), len(result))
	}
}
