package diff

import (
	"testing"
)

func makeMaskChanges() []Change {
	return []Change{
		{Key: "username", OldValue: "admin", NewValue: "root", Type: ChangeTypeModified},
		{Key: "password", OldValue: "hunter2", NewValue: "s3cr3t", Type: ChangeTypeModified},
		{Key: "api_key", OldValue: "", NewValue: "abc123", Type: ChangeTypeAdded},
		{Key: "token", OldValue: "tok_old", NewValue: "", Type: ChangeTypeRemoved},
		{Key: "host", OldValue: "localhost", NewValue: "prod.example.com", Type: ChangeTypeModified},
		{Key: "private_key", OldValue: "-----BEGIN", NewValue: "-----BEGIN NEW", Type: ChangeTypeModified},
	}
}

func TestApplyMask_Disabled(t *testing.T) {
	changes := makeMaskChanges()
	opts := DefaultMaskOptions()
	opts.Enabled = false
	result := ApplyMask(changes, opts)
	if result[1].OldValue != "hunter2" {
		t.Errorf("expected unmasked value, got %q", result[1].OldValue)
	}
}

func TestApplyMask_MasksPassword(t *testing.T) {
	changes := makeMaskChanges()
	opts := DefaultMaskOptions()
	opts.Enabled = true
	result := ApplyMask(changes, opts)
	if result[1].OldValue != "***" {
		t.Errorf("expected masked old password, got %q", result[1].OldValue)
	}
	if result[1].NewValue != "***" {
		t.Errorf("expected masked new password, got %q", result[1].NewValue)
	}
}

func TestApplyMask_MasksAPIKey(t *testing.T) {
	changes := makeMaskChanges()
	opts := DefaultMaskOptions()
	opts.Enabled = true
	result := ApplyMask(changes, opts)
	// api_key added: OldValue is empty, NewValue should be masked
	if result[2].OldValue != "" {
		t.Errorf("expected empty old value preserved, got %q", result[2].OldValue)
	}
	if result[2].NewValue != "***" {
		t.Errorf("expected masked new api_key, got %q", result[2].NewValue)
	}
}

func TestApplyMask_PreservesNonSensitive(t *testing.T) {
	changes := makeMaskChanges()
	opts := DefaultMaskOptions()
	opts.Enabled = true
	result := ApplyMask(changes, opts)
	if result[0].OldValue != "admin" {
		t.Errorf("expected username unchanged, got %q", result[0].OldValue)
	}
	if result[4].NewValue != "prod.example.com" {
		t.Errorf("expected host unchanged, got %q", result[4].NewValue)
	}
}

func TestApplyMask_CustomMaskString(t *testing.T) {
	changes := makeMaskChanges()
	opts := DefaultMaskOptions()
	opts.Enabled = true
	opts.MaskString = "[REDACTED]"
	result := ApplyMask(changes, opts)
	if result[1].NewValue != "[REDACTED]" {
		t.Errorf("expected custom mask string, got %q", result[1].NewValue)
	}
}

func TestApplyMask_EmptyChanges(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.Enabled = true
	result := ApplyMask([]Change{}, opts)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}
