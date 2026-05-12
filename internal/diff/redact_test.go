package diff

import (
	"regexp"
	"testing"
)

func makeRedactChanges() []Change {
	return []Change{
		{Key: "private_key", ChangeType: ChangeModified, OldValue: "old-pem", NewValue: "new-pem"},
		{Key: "username", ChangeType: ChangeUnchanged, OldValue: "admin", NewValue: "admin"},
		{Key: "token", ChangeType: ChangeAdded, OldValue: "", NewValue: "abc123"},
		{Key: "host", ChangeType: ChangeUnchanged, OldValue: "localhost", NewValue: "localhost"},
	}
}

func TestApplyRedact_Disabled(t *testing.T) {
	changes := makeRedactChanges()
	opts := DefaultRedactOptions()
	result := ApplyRedact(changes, opts)
	if result[0].OldValue != "old-pem" {
		t.Errorf("expected original value, got %q", result[0].OldValue)
	}
}

func TestApplyRedact_RedactsMatchingKey(t *testing.T) {
	changes := makeRedactChanges()
	re := regexp.MustCompile(`(?i)private_key`)
	opts := RedactOptions{Enabled: true, Patterns: []*regexp.Regexp{re}}
	result := ApplyRedact(changes, opts)

	if result[0].OldValue != redactedPlaceholder {
		t.Errorf("expected REDACTED for old value, got %q", result[0].OldValue)
	}
	if result[0].NewValue != redactedPlaceholder {
		t.Errorf("expected REDACTED for new value, got %q", result[0].NewValue)
	}
	if result[1].OldValue != "admin" {
		t.Errorf("non-matching key should be unchanged")
	}
}

func TestApplyRedact_EmptyValueNotRedacted(t *testing.T) {
	changes := makeRedactChanges()
	re := regexp.MustCompile(`(?i)token`)
	opts := RedactOptions{Enabled: true, Patterns: []*regexp.Regexp{re}}
	result := ApplyRedact(changes, opts)

	if result[2].OldValue != "" {
		t.Errorf("empty old value should remain empty, got %q", result[2].OldValue)
	}
	if result[2].NewValue != redactedPlaceholder {
		t.Errorf("non-empty new value should be redacted, got %q", result[2].NewValue)
	}
}

func TestParseRedactPatterns_Empty(t *testing.T) {
	names, compiled, err := ParseRedactPatterns("")
	if err != nil || names != nil || compiled != nil {
		t.Errorf("expected nil results for empty input")
	}
}

func TestParseRedactPatterns_Valid(t *testing.T) {
	names, compiled, err := ParseRedactPatterns("private_key, token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 2 || len(compiled) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(names))
	}
}

func TestParseRedactPatterns_InvalidRegex(t *testing.T) {
	_, _, err := ParseRedactPatterns("[invalid")
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}
