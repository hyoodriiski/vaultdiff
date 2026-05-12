package diff

import (
	"testing"
)

func TestParseMaskPatterns_Empty(t *testing.T) {
	opts, err := ParseMaskPatterns("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected masking disabled for empty input")
	}
}

func TestParseMaskPatterns_SinglePattern(t *testing.T) {
	opts, err := ParseMaskPatterns("secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected masking enabled")
	}
	if len(opts.Patterns) != 1 {
		t.Errorf("expected 1 pattern, got %d", len(opts.Patterns))
	}
}

func TestParseMaskPatterns_MultiplePatterns(t *testing.T) {
	opts, err := ParseMaskPatterns("password,token,api_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Patterns) != 3 {
		t.Errorf("expected 3 patterns, got %d", len(opts.Patterns))
	}
}

func TestParseMaskPatterns_InvalidRegex(t *testing.T) {
	_, err := ParseMaskPatterns("[invalid")
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestParseMaskPatterns_WhitespaceTrimmed(t *testing.T) {
	opts, err := ParseMaskPatterns(" password , token ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Patterns) != 2 {
		t.Errorf("expected 2 patterns after trim, got %d", len(opts.Patterns))
	}
}

func TestParseMaskPatterns_DefaultPatternsRetained(t *testing.T) {
	opts, err := ParseMaskPatterns("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defaults := DefaultMaskOptions()
	if len(opts.Patterns) != len(defaults.Patterns) {
		t.Errorf("expected default pattern count %d, got %d", len(defaults.Patterns), len(opts.Patterns))
	}
}
