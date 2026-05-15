package diff

import "fmt"

// SupportedPatchFormats lists valid patch output formats.
var SupportedPatchFormats = []string{"json", "shell"}

// ParsePatchFlags parses patch-related CLI flags into PatchOptions.
func ParsePatchFlags(enabled bool, outputPath, format string, dryRun bool) (PatchOptions, error) {
	opts := DefaultPatchOptions()

	if !enabled {
		return opts, nil
	}

	opts.Enabled = true
	opts.DryRun = dryRun

	if outputPath != "" {
		opts.OutputPath = outputPath
	}

	if format != "" {
		if !isSupportedPatchFormat(format) {
			return opts, fmt.Errorf("patch: unsupported format %q; supported: %v", format, SupportedPatchFormats)
		}
		opts.Format = format
	}

	return opts, nil
}

func isSupportedPatchFormat(f string) bool {
	for _, s := range SupportedPatchFormats {
		if s == f {
			return true
		}
	}
	return false
}
