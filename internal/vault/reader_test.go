package vault

import (
	"context"
	"testing"
)

func TestNormalizeKVv2Path_AlreadyNormalized(t *testing.T) {
	input := "secret/data/myapp"
	got := NormalizeKVv2Path(input)
	if got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestNormalizeKVv2Path_AddsDataSegment(t *testing.T) {
	input := "secret/myapp"
	want := "secret/data/myapp"
	got := NormalizeKVv2Path(input)
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestNormalizeKVv2Path_NoSlash(t *testing.T) {
	input := "secret"
	got := NormalizeKVv2Path(input)
	if got != input {
		t.Errorf("expected unchanged %q, got %q", input, got)
	}
}

func TestNormalizeKVv2Path_NestedPath(t *testing.T) {
	input := "kv/apps/prod/config"
	want := "kv/data/apps/prod/config"
	got := NormalizeKVv2Path(input)
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestReadSecret_KVv2Nested(t *testing.T) {
	client, mockLogical := newTestClient(t)

	mockLogical.data["secret/data/myapp"] = map[string]interface{}{
		"data": map[string]interface{}{
			"username": "admin",
			"password": "s3cr3t",
		},
	}

	result, err := client.ReadSecret(context.Background(), "secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["username"] != "admin" {
		t.Errorf("expected username=admin, got %v", result["username"])
	}
	if result["password"] != "s3cr3t" {
		t.Errorf("expected password=s3cr3t, got %v", result["password"])
	}
}

func TestReadSecret_NotFoundReader(t *testing.T) {
	client, _ := newTestClient(t)

	_, err := client.ReadSecret(context.Background(), "secret/nonexistent")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
