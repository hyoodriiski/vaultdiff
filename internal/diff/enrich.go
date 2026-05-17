package diff

import (
	"fmt"
	"strings"
)

// EnrichOptions controls metadata enrichment applied to each Change.
type EnrichOptions struct {
	Enabled        bool
	AddKeyLength   bool
	AddValueLength bool
	AddChangeID    bool
	IDPrefix       string
}

// DefaultEnrichOptions returns enrichment disabled by default.
func DefaultEnrichOptions() EnrichOptions {
	return EnrichOptions{
		Enabled:        false,
		AddKeyLength:   true,
		AddValueLength: true,
		AddChangeID:    true,
		IDPrefix:       "chg",
	}
}

// ApplyEnrich annotates each Change with optional metadata such as key/value
// lengths and a deterministic change ID derived from the key and change type.
func ApplyEnrich(changes []Change, opts EnrichOptions) []Change {
	if !opts.Enabled {
		return changes
	}

	result := make([]Change, len(changes))
	for i, c := range changes {
		if opts.AddKeyLength {
			c.Annotations = appendUniqueAnnotation(c.Annotations,
				fmt.Sprintf("key_len=%d", len(c.Key)))
		}
		if opts.AddValueLength {
			old := fmt.Sprintf("%v", c.OldValue)
			new := fmt.Sprintf("%v", c.NewValue)
			if c.Type != ChangeTypeAdded {
				c.Annotations = appendUniqueAnnotation(c.Annotations,
					fmt.Sprintf("old_len=%d", len(old)))
			}
			if c.Type != ChangeTypeRemoved {
				c.Annotations = appendUniqueAnnotation(c.Annotations,
					fmt.Sprintf("new_len=%d", len(new)))
			}
		}
		if opts.AddChangeID {
			id := buildChangeID(opts.IDPrefix, c.Key, string(c.Type))
			c.Annotations = appendUniqueAnnotation(c.Annotations, fmt.Sprintf("id=%s", id))
		}
		result[i] = c
	}
	return result
}

func buildChangeID(prefix, key, changeType string) string {
	sanitized := strings.NewReplacer("/", "_", " ", "_", ".", "_").Replace(key)
	return fmt.Sprintf("%s_%s_%s", prefix, sanitized, strings.ToLower(changeType))
}

func appendUniqueAnnotation(annotations []string, value string) []string {
	for _, a := range annotations {
		if a == value {
			return annotations
		}
	}
	return append(annotations, value)
}
