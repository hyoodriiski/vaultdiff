package diff

import (
	"testing"
)

func makeTagChanges() []Change {
	return []Change{
		{Key: "db_password", OldValue: "old", NewValue: "new", Type: Modified},
		{Key: "api_key", OldValue: "", NewValue: "abc", Type: Added},
		{Key: "host", OldValue: "a", NewValue: "a", Type: Unchanged},
		{Key: "db_host", OldValue: "x", NewValue: "y", Type: Modified},
	}
}

func TestApplyTag_Disabled(t *testing.T) {
	changes := makeTagChanges()
	opts := DefaultTagOptions()
	result := ApplyTag(changes, opts)
	for _, c := range result {
		if len(c.Annotations) != 0 {
			t.Errorf("expected no annotations when disabled, got %v for key %s", c.Annotations, c.Key)
		}
	}
}

func TestApplyTag_SinglePattern(t *testing.T) {
	changes := makeTagChanges()
	opts := TagOptions{
		Enabled: true,
		Tags:    map[string]string{"db_": "database"},
	}
	result := ApplyTag(changes, opts)

	tagged := map[string][]string{}
	for _, c := range result {
		tagged[c.Key] = c.Annotations
	}

	if !containsAnnotation(tagged["db_password"], "database") {
		t.Error("expected db_password to be tagged with 'database'")
	}
	if !containsAnnotation(tagged["db_host"], "database") {
		t.Error("expected db_host to be tagged with 'database'")
	}
	if containsAnnotation(tagged["api_key"], "database") {
		t.Error("expected api_key NOT to be tagged with 'database'")
	}
}

func TestApplyTag_MultiplePatterns(t *testing.T) {
	changes := makeTagChanges()
	opts := TagOptions{
		Enabled: true,
		Tags: map[string]string{
			"db_":      "database",
			"password": "sensitive",
		},
	}
	result := ApplyTag(changes, opts)

	for _, c := range result {
		if c.Key == "db_password" {
			if !containsAnnotation(c.Annotations, "database") {
				t.Error("expected 'database' tag on db_password")
			}
			if !containsAnnotation(c.Annotations, "sensitive") {
				t.Error("expected 'sensitive' tag on db_password")
			}
		}
	}
}

func TestApplyTag_NoDuplicateAnnotations(t *testing.T) {
	changes := []Change{
		{Key: "db_password", Type: Modified, Annotations: []string{"database"}},
	}
	opts := TagOptions{
		Enabled: true,
		Tags:    map[string]string{"db_": "database"},
	}
	result := ApplyTag(changes, opts)
	count := 0
	for _, a := range result[0].Annotations {
		if a == "database" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly one 'database' annotation, got %d", count)
	}
}

func TestParseTagFlags_Empty(t *testing.T) {
	opts, err := ParseTagFlags(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected disabled when no flags")
	}
}

func TestParseTagFlags_Valid(t *testing.T) {
	opts, err := ParseTagFlags([]string{"secret=sensitive", "db_=database"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected enabled")
	}
	if opts.Tags["secret"] != "sensitive" {
		t.Errorf("expected 'sensitive', got %q", opts.Tags["secret"])
	}
	if opts.Tags["db_"] != "database" {
		t.Errorf("expected 'database', got %q", opts.Tags["db_"])
	}
}

func containsAnnotation(annotations []string, s string) bool {
	for _, a := range annotations {
		if a == s {
			return true
		}
	}
	return false
}
