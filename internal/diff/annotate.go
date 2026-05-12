package diff

import "fmt"

// AnnotateOptions controls how change annotations are applied.
type AnnotateOptions struct {
	Enabled    bool
	ShowIndex  bool
	ShowSource bool
	SourceA    string
	SourceB    string
}

// DefaultAnnotateOptions returns sensible annotation defaults.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{
		Enabled:    false,
		ShowIndex:  false,
		ShowSource: false,
		SourceA:    "a",
		SourceB:    "b",
	}
}

// ApplyAnnotate enriches each Change with positional and source metadata.
func ApplyAnnotate(changes []Change, opts AnnotateOptions) []Change {
	if !opts.Enabled {
		return changes
	}

	annotated := make([]Change, len(changes))
	for i, c := range changes {
		if opts.ShowIndex {
			c.Key = fmt.Sprintf("%s [%d]", c.Key, i+1)
		}
		if opts.ShowSource {
			switch c.Type {
			case ChangeAdded:
				c.NewValue = fmt.Sprintf("%s (from %s)", c.NewValue, opts.SourceB)
			case ChangeRemoved:
				c.OldValue = fmt.Sprintf("%s (from %s)", c.OldValue, opts.SourceA)
			case ChangeModified:
				c.OldValue = fmt.Sprintf("%s (from %s)", c.OldValue, opts.SourceA)
				c.NewValue = fmt.Sprintf("%s (from %s)", c.NewValue, opts.SourceB)
			}
		}
		annotated[i] = c
	}
	return annotated
}
