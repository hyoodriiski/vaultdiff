package diff

import "testing"

// TestGroupAfterCompare verifies Group works correctly on real Compare output.
func TestGroupAfterCompare(t *testing.T) {
	secretA := map[string]interface{}{
		"db:host":    "localhost",
		"db:port":    "5432",
		"app:name":   "myapp",
		"app:secret": "old-secret",
	}
	secretB := map[string]interface{}{
		"db:host":    "remotehost",
		"db:port":    "5432",
		"app:name":   "myapp",
		"app:token":  "new-token",
	}

	changes := Compare(secretA, secretB)

	opts := GroupOptions{Enabled: true, GroupBy: "prefix", Separator: ":"}
	groups := Group(changes, opts)

	if len(groups) == 0 {
		t.Fatal("expected at least one group after compare")
	}

	total := 0
	for _, g := range groups {
		total += len(g.Changes)
	}
	if total != len(changes) {
		t.Errorf("expected %d total changes, got %d across groups", len(changes), total)
	}
}

// TestGroupByTypeAfterCompare verifies type-based grouping on Compare output.
func TestGroupByTypeAfterCompare(t *testing.T) {
	secretA := map[string]interface{}{
		"alpha": "a",
		"beta":  "b",
	}
	secretB := map[string]interface{}{
		"alpha": "changed",
		"gamma": "g",
	}

	changes := Compare(secretA, secretB)
	opts := GroupOptions{Enabled: true, GroupBy: "type", Separator: ":"}
	groups := Group(changes, opts)

	typesSeen := map[string]bool{}
	for _, g := range groups {
		typesSeen[g.Label] = true
		for _, c := range g.Changes {
			if string(c.Type) != g.Label {
				t.Errorf("change %q has type %q but is in group %q", c.Key, c.Type, g.Label)
			}
		}
	}

	if !typesSeen[string(ChangeTypeModified)] {
		t.Error("expected a 'modified' group")
	}
	if !typesSeen[string(ChangeTypeAdded)] {
		t.Error("expected an 'added' group")
	}
	if !typesSeen[string(ChangeTypeRemoved)] {
		t.Error("expected a 'removed' group")
	}
}
