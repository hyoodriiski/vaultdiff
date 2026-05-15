package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func makeExportChanges() []Change {
	return []Change{
		{Key: "db_host", Type: ChangeTypeModified, OldValue: "localhost", NewValue: "prod.db", Annotations: []string{"env:prod"}},
		{Key: "api_key", Type: ChangeTypeAdded, OldValue: "", NewValue: "abc123", Annotations: nil},
		{Key: "timeout", Type: ChangeTypeRemoved, OldValue: "30s", NewValue: "", Annotations: nil},
	}
}

func TestExport_Disabled(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultExportOptions()
	opts.Enabled = false
	err := Export(makeExportChanges(), opts, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output when disabled, got %q", buf.String())
	}
}

func TestExport_CSV_WithHeaders(t *testing.T) {
	var buf bytes.Buffer
	opts := ExportOptions{Enabled: true, Format: ExportFormatCSV, Headers: true}
	if err := Export(makeExportChanges(), opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines (header + 3 rows), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "key,type") {
		t.Errorf("expected header row, got %q", lines[0])
	}
}

func TestExport_CSV_NoHeaders(t *testing.T) {
	var buf bytes.Buffer
	opts := ExportOptions{Enabled: true, Format: ExportFormatCSV, Headers: false}
	if err := Export(makeExportChanges(), opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 data rows, got %d", len(lines))
	}
}

func TestExport_JSON(t *testing.T) {
	var buf bytes.Buffer
	opts := ExportOptions{Enabled: true, Format: ExportFormatJSON, Headers: true}
	if err := Export(makeExportChanges(), opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []Change
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 changes, got %d", len(out))
	}
}

func TestExport_TSV(t *testing.T) {
	var buf bytes.Buffer
	opts := ExportOptions{Enabled: true, Format: ExportFormatTSV, Headers: true}
	if err := Export(makeExportChanges(), opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\t") {
		t.Error("expected tab-separated output")
	}
}

func TestParseExportFlags_Valid(t *testing.T) {
	for _, fmt := range SupportedExportFormats {
		opts, err := ParseExportFlags(fmt, true)
		if err != nil {
			t.Errorf("format %q: unexpected error: %v", fmt, err)
		}
		if !opts.Enabled {
			t.Errorf("format %q: expected Enabled=true", fmt)
		}
	}
}

func TestParseExportFlags_Empty(t *testing.T) {
	opts, err := ParseExportFlags("", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected Enabled=false for empty format")
	}
}

func TestParseExportFlags_Invalid(t *testing.T) {
	_, err := ParseExportFlags("xml", false)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
