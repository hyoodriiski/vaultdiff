package diff

import (
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff/types"
)

// NormalizeOptions controls how values are normalized before comparison.
type NormalizeOptions struct {
	Enabled     bool
	TrimSpace   bool
	Lowercase   bool
}

// DefaultNormalizeOptions returns the default normalization settings.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		Enabled:   false,
		TrimSpace: true,
		Lowercase: false,
	}
}

// ApplyNormalize normalizes OldValue and NewValue fields on each Change
// according to the provided options. Normalization is applied before
// change-type re-evaluation so that whitespace-only differences can be
// suppressed when TrimSpace is enabled.
func ApplyNormalize(changes []types.Change, opts NormalizeOptions) []types.Change {
	if !opts.Enabled {
		return changes
	}

	result := make([]types.Change, 0, len(changes))
	for _, c := range changes {
		c.OldValue = normalizeValue(c.OldValue, opts)
		c.NewValue = normalizeValue(c.NewValue, opts)

		// Re-evaluate change type: if old and new become equal after
		// normalization, demote a Modified change to Unchanged.
		if c.ChangeType == types.ChangeTypeModified && c.OldValue == c.NewValue {
			c.ChangeType = types.ChangeTypeUnchanged
		}

		result = append(result, c)
	}
	return result
}

func normalizeValue(v string, opts NormalizeOptions) string {
	if opts.TrimSpace {
		v = strings.TrimSpace(v)
	}
	if opts.Lowercase {
		v = strings.ToLower(v)
	}
	return v
}
