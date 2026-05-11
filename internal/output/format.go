package output

import (
	"fmt"
	"io"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// Format is the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// ParseFormat converts a string to a Format, returning an error for unknowns.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf("output: unknown format %q (valid: text, json)", s)
	}
}

// Write renders the report to w in the given format.
func Write(w io.Writer, r Report, fmt Format, tf *diff.TextFormatter) error {
	switch fmt {
	case FormatJSON:
		return WriteJSON(w, r)
	case FormatText:
		return writeText(w, r, tf)
	default:
		return fmt.Errorf("output: unsupported format %q", fmt)
	}
}

func writeText(w io.Writer, r Report, tf *diff.TextFormatter) error {
	if tf == nil {
		tf = diff.NewTextFormatter()
	}
	output, err := tf.Format(r.SourcePath, r.TargetPath, r.Changes)
	if err != nil {
		return fmt.Errorf("output: rendering text: %w", err)
	}
	_, err = fmt.Fprint(w, output)
	return err
}
