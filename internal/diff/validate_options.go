package diff

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseValidateFlags parses CLI-style flag strings into ValidationOptions.
// Accepted flags:
//
//	--validate=true/false
//	--validate-reject-empty=true/false
//	--validate-reject-duplicates=true/false
//	--validate-max-key-length=N
func ParseValidateFlags(flags map[string]string) (ValidationOptions, error) {
	opts := DefaultValidateOptions()

	if v, ok := flags["validate"]; ok {
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return opts, fmt.Errorf("invalid value for --validate: %q", v)
		}
		opts.Enabled = b
	}

	if v, ok := flags["validate-reject-empty"]; ok {
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return opts, fmt.Errorf("invalid value for --validate-reject-empty: %q", v)
		}
		opts.RejectEmpty = b
	}

	if v, ok := flags["validate-reject-duplicates"]; ok {
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return opts, fmt.Errorf("invalid value for --validate-reject-duplicates: %q", v)
		}
		opts.RejectDuplicates = b
	}

	if v, ok := flags["validate-max-key-length"]; ok {
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil || n < 0 {
			return opts, fmt.Errorf("invalid value for --validate-max-key-length: %q", v)
		}
		opts.MaxKeyLength = n
	}

	return opts, nil
}
