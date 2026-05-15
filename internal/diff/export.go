package diff

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ExportFormat defines supported export formats.
type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatTSV  ExportFormat = "tsv"
)

// ExportOptions controls how changes are exported.
type ExportOptions struct {
	Enabled bool
	Format  ExportFormat
	Headers bool
}

// DefaultExportOptions returns export options with headers enabled.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Enabled: false,
		Format:  ExportFormatCSV,
		Headers: true,
	}
}

// Export writes the list of changes to w in the specified format.
func Export(changes []Change, opts ExportOptions, w io.Writer) error {
	if !opts.Enabled || len(changes) == 0 {
		return nil
	}
	switch opts.Format {
	case ExportFormatCSV:
		return exportDelimited(changes, opts.Headers, ',', w)
	case ExportFormatTSV:
		return exportDelimited(changes, opts.Headers, '\t', w)
	case ExportFormatJSON:
		return exportJSON(changes, w)
	default:
		return fmt.Errorf("unsupported export format: %q", opts.Format)
	}
}

func exportDelimited(changes []Change, headers bool, sep rune, w io.Writer) error {
	cw := csv.NewWriter(w)
	cw.Comma = sep
	if headers {
		if err := cw.Write([]string{"key", "type", "old_value", "new_value", "annotations"}); err != nil {
			return err
		}
	}
	for _, c := range changes {
		row := []string{
			c.Key,
			string(c.Type),
			c.OldValue,
			c.NewValue,
			strings.Join(c.Annotations, "|"),
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportJSON(changes []Change, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(changes)
}
