package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/output"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []string{"text", "json"}
	for _, c := range cases {
		f, err := output.ParseFormat(c)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", c, err)
		}
		if string(f) != c {
			t.Errorf("ParseFormat(%q) = %q, want %q", c, f, c)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := output.ParseFormat("yaml")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
	if !strings.Contains(err.Error(), "yaml") {
		t.Errorf("error message should mention the bad format, got: %v", err)
	}
}

func TestWrite_JSON(t *testing.T) {
	changes := []diff.Change{
		{Key: "foo", Type: diff.Added, NewValue: "bar"},
	}
	r := output.NewReport("src", "dst", changes)

	var buf bytes.Buffer
	if err := output.Write(&buf, r, output.FormatJSON, nil); err != nil {
		t.Fatalf("Write JSON error: %v", err)
	}
	if !strings.Contains(buf.String(), "\"source_path\"") {
		t.Error("JSON output missing expected field source_path")
	}
}

func TestWrite_Text(t *testing.T) {
	changes := []diff.Change{
		{Key: "key1", Type: diff.Removed, OldValue: "val"},
	}
	r := output.NewReport("src", "dst", changes)

	var buf bytes.Buffer
	if err := output.Write(&buf, r, output.FormatText, nil); err != nil {
		t.Fatalf("Write text error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty text output")
	}
}
