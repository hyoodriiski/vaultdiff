package diff

import (
	"strings"
	"testing"
)

func makeTruncateChanges() []Change {
	return []Change{
		{Key: "short", OldValue: "abc", NewValue: "def", Type: ChangeTypeModified},
		{Key: "long", OldValue: strings.Repeat("x", 120), NewValue: strings.Repeat("y", 90), Type: ChangeTypeModified},
		{Key: "exact", OldValue: strings.Repeat("z", 80), NewValue: "", Type: ChangeTypeRemoved},
	}
}

func TestApplyTruncate_Disabled(t *testing.T) {
	changes := makeTruncateChanges()
	opts := TruncateOptions{Enabled: false, MaxLength: 10, Suffix: "..."}
	result := ApplyTruncate(changes, opts)
	if len(result[1].OldValue) != 120 {
		t.Errorf("expected OldValue untouched (len 120), got len %d", len(result[1].OldValue))
	}
}

func TestApplyTruncate_ShortValueUnchanged(t *testing.T) {
	changes := makeTruncateChanges()
	opts := DefaultTruncateOptions()
	result := ApplyTruncate(changes, opts)
	if result[0].OldValue != "abc" {
		t.Errorf("expected short value unchanged, got %q", result[0].OldValue)
	}
}

func TestApplyTruncate_LongValueTruncated(t *testing.T) {
	changes := makeTruncateChanges()
	opts := DefaultTruncateOptions() // maxLen=80
	result := ApplyTruncate(changes, opts)
	got := result[1].OldValue
	if len(got) != 80 {
		t.Errorf("expected truncated length 80, got %d", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected suffix '...', got %q", got[77:])
	}
}

func TestApplyTruncate_ExactLengthUnchanged(t *testing.T) {
	changes := makeTruncateChanges()
	opts := DefaultTruncateOptions()
	result := ApplyTruncate(changes, opts)
	if len(result[2].OldValue) != 80 {
		t.Errorf("expected exact-length value unchanged (len 80), got %d", len(result[2].OldValue))
	}
}

func TestParseTruncateFlag_Default(t *testing.T) {
	opts, err := ParseTruncateFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled || opts.MaxLength != 80 {
		t.Errorf("unexpected defaults: %+v", opts)
	}
}

func TestParseTruncateFlag_Off(t *testing.T) {
	for _, flag := range []string{"off", "Off", "false", "0"} {
		opts, err := ParseTruncateFlag(flag)
		if err != nil {
			t.Fatalf("flag %q: unexpected error: %v", flag, err)
		}
		if opts.Enabled {
			t.Errorf("flag %q: expected Enabled=false", flag)
		}
	}
}

func TestParseTruncateFlag_CustomLength(t *testing.T) {
	opts, err := ParseTruncateFlag("120")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.MaxLength != 120 || !opts.Enabled {
		t.Errorf("expected MaxLength=120 Enabled=true, got %+v", opts)
	}
}

func TestParseTruncateFlag_Invalid(t *testing.T) {
	_, err := ParseTruncateFlag("banana")
	if err == nil {
		t.Error("expected error for invalid flag, got nil")
	}
}
