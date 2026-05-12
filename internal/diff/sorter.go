package diff

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering of diff changes.
type SortOrder int

const (
	SortByKey SortOrder = iota
	SortByChangeType
	SortByKeyDesc
)

// SortOptions configures how changes are sorted.
type SortOptions struct {
	Order     SortOrder
	StableKey bool // secondary sort by key when using SortByChangeType
}

// DefaultSortOptions returns the default sort configuration.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Order:     SortByKey,
		StableKey: true,
	}
}

// Sort returns a sorted copy of the provided changes slice.
func Sort(changes []Change, opts SortOptions) []Change {
	if len(changes) == 0 {
		return changes
	}

	out := make([]Change, len(changes))
	copy(out, changes)

	sort.SliceStable(out, func(i, j int) bool {
		switch opts.Order {
		case SortByChangeType:
			if out[i].Type != out[j].Type {
				return changeTypeOrder(out[i].Type) < changeTypeOrder(out[j].Type)
			}
			if opts.StableKey {
				return strings.ToLower(out[i].Key) < strings.ToLower(out[j].Key)
			}
			return false
		case SortByKeyDesc:
			return strings.ToLower(out[i].Key) > strings.ToLower(out[j].Key)
		default: // SortByKey
			return strings.ToLower(out[i].Key) < strings.ToLower(out[j].Key)
		}
	})

	return out
}

// changeTypeOrder maps ChangeType to a numeric priority for sorting.
func changeTypeOrder(ct ChangeType) int {
	switch ct {
	case ChangeTypeAdded:
		return 0
	case ChangeTypeRemoved:
		return 1
	case ChangeTypeModified:
		return 2
	case ChangeTypeUnchanged:
		return 3
	default:
		return 4
	}
}
