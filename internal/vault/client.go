package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	vc *vaultapi.Client
}

// NewClient creates a new Vault client using environment variables or provided address.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()

	if address != "" {
		cfg.Address = address
	} else if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		cfg.Address = addr
	}

	if err := cfg.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("reading vault environment: %w", err)
	}

	vc, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	if token != "" {
		vc.SetToken(token)
	} else if t := os.Getenv("VAULT_TOKEN"); t != "" {
		vc.SetToken(t)
	}

	if vc.Token() == "" {
		return nil, fmt.Errorf("vault token is required: set VAULT_TOKEN or use --token flag")
	}

	return &Client{vc: vc}, nil
}

// ReadSecret reads a KV secret at the given path and returns its data map.
func (c *Client) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := c.vc.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}

	// Support KV v2 (data is nested under "data" key)
	if data, ok := secret.Data["data"]; ok {
		if m, ok := data.(map[string]interface{}); ok {
			return m, nil
		}
	}

	return secret.Data, nil
}
