package diff

import (
	"errors"
	"fmt"
)

// ValidationResult holds the outcome of validating a set of changes.
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

// ValidationOptions controls which validation rules are applied.
type ValidationOptions struct {
	Enabled          bool
	RejectEmpty      bool
	RejectDuplicates bool
	MaxKeyLength     int
}

// DefaultValidateOptions returns sensible validation defaults.
func DefaultValidateOptions() ValidationOptions {
	return ValidationOptions{
		Enabled:          true,
		RejectEmpty:      false,
		RejectDuplicates: true,
		MaxKeyLength:     256,
	}
}

// Validate checks a slice of Changes against the provided options and returns
// a ValidationResult describing any errors or warnings found.
func Validate(changes []Change, opts ValidationOptions) (ValidationResult, error) {
	result := ValidationResult{Valid: true}

	if !opts.Enabled {
		return result, nil
	}

	seen := make(map[string]bool)

	for i, c := range changes {
		if c.Key == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("change[%d]: empty key is not allowed", i))
			result.Valid = false
			continue
		}

		if opts.MaxKeyLength > 0 && len(c.Key) > opts.MaxKeyLength {
			result.Errors = append(result.Errors,
				fmt.Sprintf("change[%d]: key %q exceeds max length %d", i, c.Key, opts.MaxKeyLength))
			result.Valid = false
		}

		if opts.RejectDuplicates {
			if seen[c.Key] {
				result.Errors = append(result.Errors,
					fmt.Sprintf("change[%d]: duplicate key %q", i, c.Key))
				result.Valid = false
			}
			seen[c.Key] = true
		}

		if opts.RejectEmpty && c.OldValue == "" && c.NewValue == "" {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("change[%d]: key %q has empty old and new values", i, c.Key))
		}
	}

	if !result.Valid {
		return result, errors.New("validation failed: one or more changes are invalid")
	}

	return result, nil
}
