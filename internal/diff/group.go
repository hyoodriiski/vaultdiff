package diff

import "sort"

// GroupOptions controls how changes are grouped.
type GroupOptions struct {
	Enabled  bool
	GroupBy  string // "type" or "prefix"
	Separator string
}

// DefaultGroupOptions returns grouping disabled by default.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		Enabled:   false,
		GroupBy:   "type",
		Separator: ":",
	}
}

// GroupedChanges holds changes organized by a group label.
type GroupedChanges struct {
	Label   string
	Changes []Change
}

// Group organizes a flat list of changes into labeled groups.
// When disabled, returns a single group with an empty label.
func Group(changes []Change, opts GroupOptions) []GroupedChanges {
	if !opts.Enabled || len(changes) == 0 {
		return []GroupedChanges{{Label: "", Changes: changes}}
	}

	switch opts.GroupBy {
	case "prefix":
		return groupByPrefix(changes, opts.Separator)
	default:
		return groupByType(changes)
	}
}

func groupByType(changes []Change) []GroupedChanges {
	buckets := map[string][]Change{}
	order := []string{}

	for _, c := range changes {
		label := string(c.Type)
		if _, exists := buckets[label]; !exists {
			order = append(order, label)
		}
		buckets[label] = append(buckets[label], c)
	}

	result := make([]GroupedChanges, 0, len(order))
	for _, label := range order {
		result = append(result, GroupedChanges{Label: label, Changes: buckets[label]})
	}
	return result
}

func groupByPrefix(changes []Change, sep string) []GroupedChanges {
	buckets := map[string][]Change{}

	for _, c := range changes {
		prefix := extractPrefix(c.Key, sep)
		buckets[prefix] = append(buckets[prefix], c)
	}

	labels := make([]string, 0, len(buckets))
	for k := range buckets {
		labels = append(labels, k)
	}
	sort.Strings(labels)

	result := make([]GroupedChanges, 0, len(labels))
	for _, label := range labels {
		result = append(result, GroupedChanges{Label: label, Changes: buckets[label]})
	}
	return result
}

func extractPrefix(key, sep string) string {
	for i, ch := range key {
		if string(ch) == sep {
			return key[:i]
		}
	}
	return key
}
