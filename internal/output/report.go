package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// Report represents a structured diff report.
type Report struct {
	GeneratedAt time.Time        `json:"generated_at"`
	SourcePath  string           `json:"source_path"`
	TargetPath  string           `json:"target_path"`
	Changes     []diff.Change    `json:"changes"`
	Summary     Summary          `json:"summary"`
}

// Summary holds counts of each change type.
type Summary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
	Unchanged int `json:"unchanged"`
	Total    int `json:"total"`
}

// NewReport builds a Report from a list of changes and path metadata.
func NewReport(sourcePath, targetPath string, changes []diff.Change) Report {
	s := Summary{Total: len(changes)}
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Modified:
			s.Modified++
		case diff.Unchanged:
			s.Unchanged++
		}
	}
	return Report{
		GeneratedAt: time.Now().UTC(),
		SourcePath:  sourcePath,
		TargetPath:  targetPath,
		Changes:     changes,
		Summary:     s,
	}
}

// WriteJSON serialises the report as indented JSON to w.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("output: encoding report as JSON: %w", err)
	}
	return nil
}
