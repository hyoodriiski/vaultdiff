package diff

import (
	"fmt"
	"io"
	"strings"
)

// Format controls the output format of a diff result.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// TextFormatter writes a human-readable diff to the given writer.
func TextFormatter(w io.Writer, result *Result) error {
	fmt.Fprintf(w, "Diff: %s → %s\n", result.PathA, result.PathB)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	if !result.HasChanges() {
		fmt.Fprintln(w, "No changes detected.")
		return nil
	}

	for _, c := range result.Changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "+ %-20s = %v\n", c.Key, c.NewValue)
		case Removed:
			fmt.Fprintf(w, "- %-20s = %v\n", c.Key, c.OldValue)
		case Modified:
			fmt.Fprintf(w, "~ %-20s : %v → %v\n", c.Key, c.OldValue, c.NewValue)
		case Unchanged:
			// omit unchanged keys in text output
		}
	}
	return nil
}
