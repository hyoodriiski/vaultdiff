package diff

import "testing"

func TestParsePatchFlags_Disabled(t *testing.T) {
	opts, err := ParsePatchFlags(false, "", "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected Enabled=false")
	}
}

func TestParsePatchFlags_DefaultFormat(t *testing.T) {
	opts, err := ParsePatchFlags(true, "", "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != "json" {
		t.Errorf("expected default format 'json', got %q", opts.Format)
	}
}

func TestParsePatchFlags_ValidFormat(t *testing.T) {
	for _, f := range SupportedPatchFormats {
		_, err := ParsePatchFlags(true, "", f, false)
		if err != nil {
			t.Errorf("format %q should be valid, got error: %v", f, err)
		}
	}
}

func TestParsePatchFlags_InvalidFormat(t *testing.T) {
	_, err := ParsePatchFlags(true, "", "xml", false)
	if err == nil {
		t.Error("expected error for unsupported format 'xml'")
	}
}

func TestParsePatchFlags_OutputPath(t *testing.T) {
	opts, err := ParsePatchFlags(true, "/tmp/out.json", "json", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.OutputPath != "/tmp/out.json" {
		t.Errorf("expected OutputPath '/tmp/out.json', got %q", opts.OutputPath)
	}
}

func TestParsePatchFlags_DryRun(t *testing.T) {
	opts, err := ParsePatchFlags(true, "", "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.DryRun {
		t.Error("expected DryRun=true")
	}
}
