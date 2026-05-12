package diff

import (
	"fmt"
	"strings"
)

// SupportedSortOrders lists valid sort order flag values.
var SupportedSortOrders = []string{"key", "key-desc", "type"}

// ParseSortFlag converts a string flag value into a SortOptions struct.
func ParseSortFlag(s string) (SortOptions, error) {
	opts := DefaultSortOptions()

	if s == "" {
		return opts, nil
	}

	switch strings.ToLower(strings.TrimSpace(s)) {
	case "key":
		opts.Order = SortByKey
	case "key-desc":
		opts.Order = SortByKeyDesc
	case "type":
		opts.Order = SortByChangeType
		opts.StableKey = true
	default:
		return opts, fmt.Errorf(
			"unsupported sort order %q: must be one of %s",
			s,
			strings.Join(SupportedSortOrders, ", "),
		)
	}

	return opts, nil
}
