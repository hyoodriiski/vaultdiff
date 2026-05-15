package diff

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested map values are flattened into dot-notation keys.
type FlattenOptions struct {
	Enabled   bool
	Separator string
	MaxDepth  int
}

// DefaultFlattenOptions returns sensible defaults for flattening.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Enabled:   false,
		Separator: ".",
		MaxDepth:  10,
	}
}

// ApplyFlatten expands any Change whose OldValue or NewValue contains
// JSON-like nested map notation (key=map[...]) into individual dot-keyed
// Changes. Non-nested changes are passed through unchanged.
func ApplyFlatten(changes []Change, opts FlattenOptions) []Change {
	if !opts.Enabled {
		return changes
	}
	if opts.Separator == "" {
		opts.Separator = "."
	}

	result := make([]Change, 0, len(changes))
	for _, c := range changes {
		flattened := flattenChange(c, opts)
		result = append(result, flattened...)
	}
	return result
}

func flattenChange(c Change, opts FlattenOptions) []Change {
	oldMap := tryParseNestedValue(c.OldValue)
	newMap := tryParseNestedValue(c.NewValue)

	if oldMap == nil && newMap == nil {
		return []Change{c}
	}

	keys := unionStringKeys(oldMap, newMap)
	sort.Strings(keys)

	result := make([]Change, 0, len(keys))
	for _, k := range keys {
		subKey := joinKey(c.Key, k, opts.Separator)
		oldVal := mapGet(oldMap, k)
		newVal := mapGet(newMap, k)

		sub := Change{
			Key:      subKey,
			OldValue: oldVal,
			NewValue: newVal,
		}
		sub.Type = inferChangeType(oldVal, newVal)
		result = append(result, sub)
	}
	return result
}

func tryParseNestedValue(v string) map[string]string {
	v = strings.TrimSpace(v)
	if !strings.HasPrefix(v, "map[") {
		return nil
	}
	inner := strings.TrimPrefix(v, "map[")
	inner = strings.TrimSuffix(inner, "]")
	pairs := strings.Fields(inner)
	result := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, ":", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func unionStringKeys(a, b map[string]string) []string {
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

func mapGet(m map[string]string, key string) string {
	if m == nil {
		return ""
	}
	return m[key]
}

func joinKey(parent, child, sep string) string {
	if parent == "" {
		return child
	}
	return fmt.Sprintf("%s%s%s", parent, sep, child)
}

func inferChangeType(oldVal, newVal string) ChangeType {
	switch {
	case oldVal == "" && newVal != "":
		return ChangeTypeAdded
	case oldVal != "" && newVal == "":
		return ChangeTypeRemoved
	case oldVal != newVal:
		return ChangeTypeModified
	default:
		return ChangeTypeUnchanged
	}
}
