package diff

import (
	"fmt"
	"strings"
)

// SupportedExportFormats lists all valid export format strings.
var SupportedExportFormats = []string{
	string(ExportFormatCSV),
	string(ExportFormatJSON),
	string(ExportFormatTSV),
}

// ParseExportFlags parses the --export-format flag value into ExportOptions.
// An empty format string disables export.
func ParseExportFlags(format string, headers bool) (ExportOptions, error) {
	opts := DefaultExportOptions()
	opts.Headers = headers

	format = strings.TrimSpace(strings.ToLower(format))
	if format == "" {
		opts.Enabled = false
		return opts, nil
	}

	switch ExportFormat(format) {
	case ExportFormatCSV, ExportFormatJSON, ExportFormatTSV:
		opts.Enabled = true
		opts.Format = ExportFormat(format)
		return opts, nil
	default:
		return opts, fmt.Errorf(
			"invalid export format %q: must be one of [%s]",
			format,
			strings.Join(SupportedExportFormats, ", "),
		)
	}
}

// IsValidExportFormat reports whether the given format string is a supported
// export format. The check is case-insensitive.
func IsValidExportFormat(format string) bool {
	format = strings.TrimSpace(strings.ToLower(format))
	switch ExportFormat(format) {
	case ExportFormatCSV, ExportFormatJSON, ExportFormatTSV:
		return true
	default:
		return false
	}
}
