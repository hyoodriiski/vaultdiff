package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// CheckpointOptions controls checkpoint save/load behaviour.
type CheckpointOptions struct {
	Enabled  bool
	Path     string
	AutoLoad bool
}

// DefaultCheckpointOptions returns safe defaults (disabled).
func DefaultCheckpointOptions() CheckpointOptions {
	return CheckpointOptions{
		Enabled:  false,
		Path:     "",
		AutoLoad: false,
	}
}

// checkpointFile is the on-disk representation of a saved checkpoint.
type checkpointFile struct {
	SavedAt time.Time `json:"saved_at"`
	Changes []Change  `json:"changes"`
}

// SaveCheckpoint writes the current changes to a checkpoint file.
func SaveCheckpoint(opts CheckpointOptions, changes []Change) error {
	if !opts.Enabled || opts.Path == "" {
		return nil
	}

	cf := checkpointFile{
		SavedAt: time.Now().UTC(),
		Changes: changes,
	}

	data, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}

	if err := os.WriteFile(opts.Path, data, 0o600); err != nil {
		return fmt.Errorf("checkpoint: write %s: %w", opts.Path, err)
	}

	return nil
}

// LoadCheckpoint reads a previously saved checkpoint file.
// Returns nil slice when the file does not exist and opts.AutoLoad is false.
func LoadCheckpoint(opts CheckpointOptions) ([]Change, time.Time, error) {
	if !opts.Enabled || opts.Path == "" {
		return nil, time.Time{}, nil
	}

	data, err := os.ReadFile(opts.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, time.Time{}, nil
		}
		return nil, time.Time{}, fmt.Errorf("checkpoint: read %s: %w", opts.Path, err)
	}

	var cf checkpointFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, time.Time{}, fmt.Errorf("checkpoint: unmarshal: %w", err)
	}

	return cf.Changes, cf.SavedAt, nil
}
