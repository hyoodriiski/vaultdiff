package vault

import (
	"context"
	"fmt"
	"strings"
)

// SecretData represents the key-value pairs stored at a Vault path.
type SecretData map[string]interface{}

// ReadSecret reads a secret from the given path using the KV v1 or v2 engine.
// For KV v2, the path is automatically rewritten to include /data/.
func (c *Client) ReadSecret(ctx context.Context, path string) (SecretData, error) {
	secret, err := c.logical.ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}

	// KV v2 stores data nested under a "data" key
	if data, ok := secret.Data["data"]; ok {
		if nested, ok := data.(map[string]interface{}); ok {
			return SecretData(nested), nil
		}
	}

	return SecretData(secret.Data), nil
}

// NormalizeKVv2Path rewrites a KV v2 path to include the /data/ segment
// if it is not already present. For example:
//   "secret/myapp" -> "secret/data/myapp"
func NormalizeKVv2Path(path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return path
	}
	mount := parts[0]
	rest := parts[1]
	if strings.HasPrefix(rest, "data/") {
		return path
	}
	return mount + "/data/" + rest
}
