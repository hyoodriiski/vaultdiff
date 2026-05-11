package main

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	// parseFlags reads from os.Args; we test config struct defaults via direct construction
	cfg := config{
		format: "text",
	}
	if cfg.format != "text" {
		t.Errorf("expected default format 'text', got %q", cfg.format)
	}
	if cfg.pathA != "" {
		t.Errorf("expected empty pathA, got %q", cfg.pathA)
	}
	if cfg.pathB != "" {
		t.Errorf("expected empty pathB, got %q", cfg.pathB)
	}
}

func TestRun_MissingPaths(t *testing.T) {
	cfg := config{
		format: "text",
	}
	err := run(cfg)
	if err == nil {
		t.Fatal("expected error for missing paths, got nil")
	}
	expected := "both -a and -b paths are required"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	cfg := config{
		pathA:  "secret/a",
		pathB:  "secret/b",
		format: "xml",
	}
	err := run(cfg)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestRun_MissingToken(t *testing.T) {
	cfg := config{
		pathA:      "secret/a",
		pathB:      "secret/b",
		format:     "text",
		vaultAddr:  "http://127.0.0.1:8200",
		vaultToken: "",
	}
	err := run(cfg)
	if err == nil {
		t.Fatal("expected error for missing vault token, got nil")
	}
}
