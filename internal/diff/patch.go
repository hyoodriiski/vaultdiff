package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PatchOptions controls patch generation behavior.
type PatchOptions struct {
	Enabled   bool
	OutputPath string
	Format    string // "json" or "shell"
	DryRun    bool
}

// DefaultPatchOptions returns safe defaults.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{
		Enabled:   false,
		OutputPath: "",
		Format:    "json",
		DryRun:    false,
	}
}

// PatchEntry represents a single reversible change.
type PatchEntry struct {
	Key      string `json:"key"`
	Op       string `json:"op"` // "set", "delete"
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

// Patch holds a full set of patch entries with metadata.
type Patch struct {
	GeneratedAt string       `json:"generated_at"`
	Entries     []PatchEntry `json:"entries"`
}

// GeneratePatch builds a Patch from a slice of Change.
func GeneratePatch(changes []Change, opts PatchOptions) (*Patch, error) {
	if !opts.Enabled {
		return nil, nil
	}

	patch := &Patch{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Entries:     make([]PatchEntry, 0, len(changes)),
	}

	for _, c := range changes {
		switch c.Type {
		case ChangeTypeAdded:
			patch.Entries = append(patch.Entries, PatchEntry{
				Key:      c.Key,
				Op:       "set",
				NewValue: fmt.Sprintf("%v", c.NewValue),
			})
		case ChangeTypeRemoved:
			patch.Entries = append(patch.Entries, PatchEntry{
				Key:      c.Key,
				Op:       "delete",
				OldValue: fmt.Sprintf("%v", c.OldValue),
			})
		case ChangeTypeModified:
			patch.Entries = append(patch.Entries, PatchEntry{
				Key:      c.Key,
				Op:       "set",
				OldValue: fmt.Sprintf("%v", c.OldValue),
				NewValue: fmt.Sprintf("%v", c.NewValue),
			})
		}
	}

	if !opts.DryRun && opts.OutputPath != "" {
		if err := writePatch(patch, opts.OutputPath); err != nil {
			return patch, err
		}
	}

	return patch, nil
}

func writePatch(patch *Patch, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("patch: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(patch)
}
