package diff

import (
	"encoding/json"
	"fmt"
	"os"
)

// BaselineOptions controls baseline snapshot behavior.
type BaselineOptions struct {
	Enabled  bool
	FilePath string
}

// DefaultBaselineOptions returns baseline options with saving disabled.
func DefaultBaselineOptions() BaselineOptions {
	return BaselineOptions{
		Enabled:  false,
		FilePath: "",
	}
}

// baselineSnapshot represents a persisted set of secret key/value pairs.
type baselineSnapshot struct {
	Data map[string]string `json:"data"`
}

// SaveBaseline writes the current "left" values of all changes to a JSON file
// so they can be used as a future comparison baseline.
func SaveBaseline(changes []Change, opts BaselineOptions) error {
	if !opts.Enabled || opts.FilePath == "" {
		return nil
	}

	snap := baselineSnapshot{
		Data: make(map[string]string, len(changes)),
	}

	for _, c := range changes {
		switch c.Type {
		case ChangeTypeAdded:
			snap.Data[c.Key] = c.NewValue
		case ChangeTypeRemoved:
			snap.Data[c.Key] = c.OldValue
		case ChangeTypeModified:
			snap.Data[c.Key] = c.NewValue
		case ChangeTypeUnchanged:
			snap.Data[c.Key] = c.OldValue
		}
	}

	f, err := os.Create(opts.FilePath)
	if err != nil {
		return fmt.Errorf("baseline: create file %q: %w", opts.FilePath, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("baseline: encode snapshot: %w", err)
	}
	return nil
}

// LoadBaseline reads a previously saved baseline snapshot and returns it as a
// flat map of key → value suitable for use as the "left" side of a Compare call.
func LoadBaseline(filePath string) (map[string]interface{}, error) {
	if filePath == "" {
		return nil, fmt.Errorf("baseline: file path is empty")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("baseline: open file %q: %w", filePath, err)
	}
	defer f.Close()

	var snap baselineSnapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("baseline: decode snapshot: %w", err)
	}

	out := make(map[string]interface{}, len(snap.Data))
	for k, v := range snap.Data {
		out[k] = v
	}
	return out, nil
}
