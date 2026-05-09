# vaultdiff

> CLI tool to diff two HashiCorp Vault secret paths and output structured change reports.

---

## Installation

```bash
go install github.com/youruser/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

Ensure your Vault environment variables are set (`VAULT_ADDR`, `VAULT_TOKEN`), then run:

```bash
vaultdiff secret/app/production secret/app/staging
```

### Example Output

```
~ db_password   [changed]
+ new_api_key   [added]
- old_feature   [removed]
= log_level     [unchanged]
```

### Flags

| Flag         | Description                              | Default  |
|--------------|------------------------------------------|----------|
| `--format`   | Output format: `text`, `json`, `yaml`   | `text`   |
| `--show-all` | Include unchanged keys in output         | `false`  |
| `--mount`    | Vault KV mount path                      | `secret` |

```bash
# Output as JSON
vaultdiff --format=json secret/app/v1 secret/app/v2

# Include unchanged keys
vaultdiff --show-all secret/app/v1 secret/app/v2
```

---

## Requirements

- Go 1.21+
- HashiCorp Vault (KV v1 or v2)

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 youruser