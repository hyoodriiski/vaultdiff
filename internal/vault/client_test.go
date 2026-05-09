package vault

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewClient_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	os.Unsetenv("VAULT_ADDR")

	_, err := NewClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
}

func TestNewClient_WithToken(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestReadSecret_KVv1(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"foo":"bar","baz":"qux"}}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := client.ReadSecret("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secret: %v", err)
	}

	if data["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %v", data["foo"])
	}
}

func TestReadSecret_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.ReadSecret("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
