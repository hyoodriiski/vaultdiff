package diff

import "fmt"

// ParseCheckpointFlags builds CheckpointOptions from raw CLI values.
func ParseCheckpointFlags(enabled bool, path string, autoLoad bool) (CheckpointOptions, error) {
	if !enabled {
		return DefaultCheckpointOptions(), nil
	}

	if path == "" {
		return CheckpointOptions{}, fmt.Errorf("checkpoint: --checkpoint-path must be set when checkpoints are enabled")
	}

	return CheckpointOptions{
		Enabled:  true,
		Path:     path,
		AutoLoad: autoLoad,
	}, nil
}
