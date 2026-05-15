package diff

import "strings"

// RenameOptions controls key rename detection in a diff.
type RenameOptions struct {
	Enabled    bool
	Similarity float64 // 0.0–1.0, minimum value similarity to consider a rename
}

// DefaultRenameOptions returns conservative rename detection settings.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		Enabled:    false,
		Similarity: 0.8,
	}
}

// ApplyRename detects removed+added pairs that look like renames and annotates
// them with a "renamed-from" / "renamed-to" tag in their metadata.
func ApplyRename(changes []Change, opts RenameOptions) []Change {
	if !opts.Enabled {
		return changes
	}

	removed := make([]int, 0)
	added := make([]int, 0)
	for i, c := range changes {
		switch c.Type {
		case ChangeTypeRemoved:
			removed = append(removed, i)
		case ChangeTypeAdded:
			added = append(added, i)
		}
	}

	used := make(map[int]bool)
	result := make([]Change, len(changes))
	copy(result, changes)

	for _, ri := range removed {
		bestScore := -1.0
		bestAi := -1
		for _, ai := range added {
			if used[ai] {
				continue
			}
			score := valueSimilarity(changes[ri].OldValue, changes[ai].NewValue)
			if score >= opts.Similarity && score > bestScore {
				bestScore = score
				bestAi = ai
			}
		}
		if bestAi >= 0 {
			used[bestAi] = true
			if result[ri].Meta == nil {
				result[ri].Meta = map[string]string{}
			}
			if result[bestAi].Meta == nil {
				result[bestAi].Meta = map[string]string{}
			}
			result[ri].Meta["renamed-to"] = changes[bestAi].Key
			result[bestAi].Meta["renamed-from"] = changes[ri].Key
		}
	}
	return result
}

// valueSimilarity returns a rough similarity score between two strings (0–1).
func valueSimilarity(a, b string) float64 {
	if a == b {
		return 1.0
	}
	if a == "" || b == "" {
		return 0.0
	}
	shorter, longer := a, b
	if len(a) > len(b) {
		shorter, longer = b, a
	}
	if strings.Contains(longer, shorter) {
		return float64(len(shorter)) / float64(len(longer))
	}
	common := 0
	for i := 0; i < len(shorter); i++ {
		if i < len(longer) && shorter[i] == longer[i] {
			common++
		}
	}
	return float64(common) / float64(len(longer))
}
