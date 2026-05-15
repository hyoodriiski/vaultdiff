package diff

import (
	"fmt"
	"strings"
)

// HighlightOptions controls inline character-level diff highlighting.
type HighlightOptions struct {
	Enabled bool
	Prefix  string
	Suffix  string
}

// DefaultHighlightOptions returns highlight options with ANSI bold markers.
func DefaultHighlightOptions() HighlightOptions {
	return HighlightOptions{
		Enabled: false,
		Prefix:  "[[",
		Suffix:  "]]",
	}
}

// ApplyHighlight annotates Modified changes with inline character-level diff
// markers showing which characters changed between OldValue and NewValue.
func ApplyHighlight(changes []Change, opts HighlightOptions) []Change {
	if !opts.Enabled {
		return changes
	}

	result := make([]Change, len(changes))
	copy(result, changes)

	for i, c := range result {
		if c.Type != ChangeTypeModified {
			continue
		}
		oldH, newH := highlightDiff(c.OldValue, c.NewValue, opts.Prefix, opts.Suffix)
		result[i].OldValue = oldH
		result[i].NewValue = newH
	}
	return result
}

// highlightDiff marks differing characters between two strings.
func highlightDiff(old, new_, prefix, suffix string) (string, string) {
	oldRunes := []rune(old)
	newRunes := []rune(new_)

	if old == new_ {
		return old, new_
	}

	// Find common prefix length.
	prefixLen := 0
	for prefixLen < len(oldRunes) && prefixLen < len(newRunes) && oldRunes[prefixLen] == newRunes[prefixLen] {
		prefixLen++
	}

	// Find common suffix length.
	suffixLen := 0
	for suffixLen < len(oldRunes)-prefixLen && suffixLen < len(newRunes)-prefixLen &&
		oldRunes[len(oldRunes)-1-suffixLen] == newRunes[len(newRunes)-1-suffixLen] {
		suffixLen++
	}

	oldMid := string(oldRunes[prefixLen : len(oldRunes)-suffixLen])
	newMid := string(newRunes[prefixLen : len(newRunes)-suffixLen])

	commonPre := string(oldRunes[:prefixLen])
	commonSuf := string(oldRunes[len(oldRunes)-suffixLen:])

	oldHighlighted := fmt.Sprintf("%s%s%s%s%s", commonPre, prefix, oldMid, suffix, commonSuf)
	newHighlighted := fmt.Sprintf("%s%s%s%s%s", commonPre, prefix, newMid, suffix, commonSuf)

	// Clean up empty highlight markers.
	empty := prefix + suffix
	oldHighlighted = strings.ReplaceAll(oldHighlighted, empty, "")
	newHighlighted = strings.ReplaceAll(newHighlighted, empty, "")

	return oldHighlighted, newHighlighted
}
