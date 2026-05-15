package diff

import (
	"testing"
)

func makeLimitChanges(n int) []Change {
	changes := make([]Change, n)
	for i := 0; i < n; i++ {
		changes[i] = Change{
			Key:      fmt.Sprintf("key%d", i),
			Type:     ChangeTypeAdded,
			NewValue: "val",
		}
	}
	return changes
}

func TestApplyLimit_Disabled(t *testing.T) {
	changes := makeLimitChanges(10)
	opts := DefaultLimitOptions()
	out, o := ApplyLimit(changes, opts)
	if len(out) != 10 {
		t.Errorf("expected 10 changes, got %d", len(out))
	}
	if o.Truncated {
		t.Error("expected Truncated=false when disabled")
	}
}

func TestApplyLimit_BelowMax(t *testing.T) {
	changes := makeLimitChanges(5)
	opts := LimitOptions{Enabled: true, MaxItems: 10}
	out, o := ApplyLimit(changes, opts)
	if len(out) != 5 {
		t.Errorf("expected 5 changes, got %d", len(out))
	}
	if o.Truncated {
		t.Error("expected Truncated=false when under limit")
	}
}

func TestApplyLimit_ExceedsMax(t *testing.T) {
	changes := makeLimitChanges(20)
	opts := LimitOptions{Enabled: true, MaxItems: 7}
	out, o := ApplyLimit(changes, opts)
	if len(out) != 7 {
		t.Errorf("expected 7 changes, got %d", len(out))
	}
	if !o.Truncated {
		t.Error("expected Truncated=true when over limit")
	}
}

func TestApplyLimit_ZeroMax(t *testing.T) {
	changes := makeLimitChanges(5)
	opts := LimitOptions{Enabled: true, MaxItems: 0}
	out, _ := ApplyLimit(changes, opts)
	if len(out) != 5 {
		t.Errorf("expected 5 changes (zero max disables), got %d", len(out))
	}
}

func TestParseLimitFlag_Empty(t *testing.T) {
	opts, err := ParseLimitFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected disabled for empty string")
	}
}

func TestParseLimitFlag_Valid(t *testing.T) {
	opts, err := ParseLimitFlag("25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected enabled")
	}
	if opts.MaxItems != 25 {
		t.Errorf("expected MaxItems=25, got %d", opts.MaxItems)
	}
}

func TestParseLimitFlag_Invalid(t *testing.T) {
	_, err := ParseLimitFlag("abc")
	if err == nil {
		t.Error("expected error for non-numeric value")
	}
}

func TestParseLimitFlag_Negative(t *testing.T) {
	_, err := ParseLimitFlag("-5")
	if err == nil {
		t.Error("expected error for negative value")
	}
}

func TestLimitSummary_NotTruncated(t *testing.T) {
	opts := LimitOptions{Enabled: true, MaxItems: 10, Truncated: false}
	if s := LimitSummary(opts, 5); s != "" {
		t.Errorf("expected empty summary, got %q", s)
	}
}

func TestLimitSummary_Truncated(t *testing.T) {
	opts := LimitOptions{Enabled: true, MaxItems: 10, Truncated: true}
	s := LimitSummary(opts, 50)
	if s == "" {
		t.Error("expected non-empty summary when truncated")
	}
}
