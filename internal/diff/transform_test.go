package diff

import "testing"

func makeTransformChanges() []Change {
	return []Change{
		{Key: "host", OldValue: "  localhost  ", NewValue: "  127.0.0.1  ", Type: ChangeTypeModified},
		{Key: "env", OldValue: "Production", NewValue: "Production", Type: ChangeTypeUnchanged},
		{Key: "token", OldValue: "", NewValue: "Bearer abc123", Type: ChangeTypeAdded},
		{Key: "old_key", OldValue: "prefix_value", NewValue: "", Type: ChangeTypeRemoved},
	}
}

func TestApplyTransform_Disabled(t *testing.T) {
	changes := makeTransformChanges()
	opts := DefaultTransformOptions()
	result := ApplyTransform(changes, opts)
	if result[0].OldValue != "  localhost  " {
		t.Errorf("expected unchanged value, got %q", result[0].OldValue)
	}
}

func TestApplyTransform_Uppercase(t *testing.T) {
	changes := makeTransformChanges()
	opts := TransformOptions{Enabled: true, Uppercase: true}
	result := ApplyTransform(changes, opts)
	if result[1].OldValue != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", result[1].OldValue)
	}
}

func TestApplyTransform_Lowercase(t *testing.T) {
	changes := makeTransformChanges()
	opts := TransformOptions{Enabled: true, Lowercase: true}
	result := ApplyTransform(changes, opts)
	if result[1].NewValue != "production" {
		t.Errorf("expected production, got %q", result[1].NewValue)
	}
}

func TestApplyTransform_TrimPrefix(t *testing.T) {
	changes := makeTransformChanges()
	opts := TransformOptions{Enabled: true, TrimPrefix: "prefix_"}
	result := ApplyTransform(changes, opts)
	if result[3].OldValue != "value" {
		t.Errorf("expected 'value', got %q", result[3].OldValue)
	}
}

func TestApplyTransform_TrimSuffix(t *testing.T) {
	changes := makeTransformChanges()
	opts := TransformOptions{Enabled: true, TrimSuffix: "123"}
	result := ApplyTransform(changes, opts)
	if result[2].NewValue != "Bearer abc" {
		t.Errorf("expected 'Bearer abc', got %q", result[2].NewValue)
	}
}

func TestApplyTransform_EmptyValueUnchanged(t *testing.T) {
	changes := makeTransformChanges()
	opts := TransformOptions{Enabled: true, Uppercase: true}
	result := ApplyTransform(changes, opts)
	if result[2].OldValue != "" {
		t.Errorf("expected empty string, got %q", result[2].OldValue)
	}
}

func TestParseTransformFlags_Empty(t *testing.T) {
	opts := ParseTransformFlags(nil)
	if opts.Enabled {
		t.Error("expected Enabled=false for nil flags")
	}
}

func TestParseTransformFlags_Uppercase(t *testing.T) {
	opts := ParseTransformFlags([]string{"uppercase"})
	if !opts.Enabled || !opts.Uppercase {
		t.Error("expected Enabled=true and Uppercase=true")
	}
}

func TestParseTransformFlags_TrimPrefix(t *testing.T) {
	opts := ParseTransformFlags([]string{"trim-prefix=Bearer "})
	if opts.TrimPrefix != "Bearer " {
		t.Errorf("expected 'Bearer ', got %q", opts.TrimPrefix)
	}
}
