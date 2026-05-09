package diff

import "sort"

// ChangeType represents the type of change detected between two secrets.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level difference between two secret maps.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue interface{}
	NewValue interface{}
}

// Result holds the full diff result between two secret paths.
type Result struct {
	PathA   string
	PathB   string
	Changes []Change
}

// HasChanges returns true if there are any non-unchanged entries.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Compare computes the diff between two secret data maps.
func Compare(pathA, pathB string, secretA, secretB map[string]interface{}) *Result {
	result := &Result{
		PathA: pathA,
		PathB: pathB,
	}

	keys := unionKeys(secretA, secretB)
	sort.Strings(keys)

	for _, key := range keys {
		valA, inA := secretA[key]
		valB, inB := secretB[key]

		var change Change
		change.Key = key

		switch {
		case inA && !inB:
			change.Type = Removed
			change.OldValue = valA
		case !inA && inB:
			change.Type = Added
			change.NewValue = valB
		case formatValue(valA) != formatValue(valB):
			change.Type = Modified
			change.OldValue = valA
			change.NewValue = valB
		default:
			change.Type = Unchanged
			change.OldValue = valA
			change.NewValue = valB
		}

		result.Changes = append(result.Changes, change)
	}

	return result
}

func unionKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
