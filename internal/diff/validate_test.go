package diff

import (
	"testing"
)

func makeValidateChanges() []Change {
	return []Change{
		{Key: "db_password", OldValue: "secret", NewValue: "newsecret", Type: ChangeTypeModified},
		{Key: "api_key", OldValue: "", NewValue: "abc123", Type: ChangeTypeAdded},
		{Key: "old_token", OldValue: "tok", NewValue: "", Type: ChangeTypeRemoved},
	}
}

func TestValidate_Disabled(t *testing.T) {
	opts := DefaultValidateOptions()
	opts.Enabled = false
	changes := []Change{{Key: "", Type: ChangeTypeAdded}}
	result, err := Validate(changes, opts)
	if err != nil {
		t.Fatalf("expected no error when disabled, got %v", err)
	}
	if !result.Valid {
		t.Error("expected Valid=true when disabled")
	}
}

func TestValidate_ValidChanges(t *testing.T) {
	opts := DefaultValidateOptions()
	changes := makeValidateChanges()
	result, err := Validate(changes, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidate_EmptyKeyRejected(t *testing.T) {
	opts := DefaultValidateOptions()
	changes := []Change{{Key: "", Type: ChangeTypeAdded, NewValue: "val"}}
	result, err := Validate(changes, opts)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
	if result.Valid {
		t.Error("expected Valid=false")
	}
	if len(result.Errors) == 0 {
		t.Error("expected at least one error message")
	}
}

func TestValidate_DuplicateKey(t *testing.T) {
	opts := DefaultValidateOptions()
	opts.RejectDuplicates = true
	changes := []Change{
		{Key: "foo", Type: ChangeTypeAdded, NewValue: "a"},
		{Key: "foo", Type: ChangeTypeModified, OldValue: "a", NewValue: "b"},
	}
	result, err := Validate(changes, opts)
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
	if result.Valid {
		t.Error("expected Valid=false for duplicate key")
	}
}

func TestValidate_MaxKeyLength(t *testing.T) {
	opts := DefaultValidateOptions()
	opts.MaxKeyLength = 5
	changes := []Change{{Key: "toolongkey", Type: ChangeTypeAdded, NewValue: "v"}}
	result, err := Validate(changes, opts)
	if err == nil {
		t.Fatal("expected error for key exceeding max length")
	}
	if result.Valid {
		t.Error("expected Valid=false")
	}
}

func TestValidate_RejectEmptyWarns(t *testing.T) {
	opts := DefaultValidateOptions()
	opts.RejectEmpty = true
	changes := []Change{{Key: "ghost", Type: ChangeTypeUnchanged, OldValue: "", NewValue: ""}}
	result, err := Validate(changes, opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Valid {
		t.Error("expected Valid=true (warnings only)")
	}
	if len(result.Warnings) == 0 {
		t.Error("expected at least one warning for empty values")
	}
}

func TestParseValidateFlags_Valid(t *testing.T) {
	flags := map[string]string{
		"validate":                   "true",
		"validate-reject-empty":      "true",
		"validate-reject-duplicates": "false",
		"validate-max-key-length":    "128",
	}
	opts, err := ParseValidateFlags(flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected Enabled=true")
	}
	if !opts.RejectEmpty {
		t.Error("expected RejectEmpty=true")
	}
	if opts.RejectDuplicates {
		t.Error("expected RejectDuplicates=false")
	}
	if opts.MaxKeyLength != 128 {
		t.Errorf("expected MaxKeyLength=128, got %d", opts.MaxKeyLength)
	}
}

func TestParseValidateFlags_InvalidBool(t *testing.T) {
	flags := map[string]string{"validate": "notabool"}
	_, err := ParseValidateFlags(flags)
	if err == nil {
		t.Fatal("expected error for invalid bool")
	}
}

func TestParseValidateFlags_InvalidMaxKeyLength(t *testing.T) {
	flags := map[string]string{"validate-max-key-length": "-5"}
	_, err := ParseValidateFlags(flags)
	if err == nil {
		t.Fatal("expected error for negative max key length")
	}
}
