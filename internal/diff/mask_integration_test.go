package diff

import (
	"testing"
)

// TestMaskAfterCompare verifies masking integrates correctly with Compare output.
func TestMaskAfterCompare(t *testing.T) {
	old := map[string]interface{}{
		"host":     "localhost",
		"password": "old_pass",
		"token":    "old_token",
	}
	new := map[string]interface{}{
		"host":     "prod.example.com",
		"password": "new_pass",
		"token":    "old_token",
	}
	changes := Compare(old, new)
	opts := DefaultMaskOptions()
	opts.Enabled = true
	masked := ApplyMask(changes, opts)

	for _, c := range masked {
		switch c.Key {
		case "host":
			if c.NewValue == "***" {
				t.Error("host should not be masked")
			}
		case "password":
			if c.OldValue != "***" || c.NewValue != "***" {
				t.Errorf("password values should be masked, got old=%q new=%q", c.OldValue, c.NewValue)
			}
		case "token":
			if c.OldValue != "***" {
				t.Errorf("token old value should be masked, got %q", c.OldValue)
			}
		}
	}
}

// TestMaskWithCustomPattern verifies a user-supplied pattern masks correctly.
func TestMaskWithCustomPattern(t *testing.T) {
	old := map[string]interface{}{
		"db_pass": "hunter2",
		"db_user": "admin",
	}
	new := map[string]interface{}{
		"db_pass": "n3wpass",
		"db_user": "superadmin",
	}
	changes := Compare(old, new)
	opts, err := ParseMaskPatterns("db_pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	masked := ApplyMask(changes, opts)

	for _, c := range masked {
		switch c.Key {
		case "db_pass":
			if c.OldValue != "***" || c.NewValue != "***" {
				t.Errorf("db_pass should be masked, got old=%q new=%q", c.OldValue, c.NewValue)
			}
		case "db_user":
			if c.OldValue == "***" || c.NewValue == "***" {
				t.Error("db_user should not be masked")
			}
		}
	}
}
