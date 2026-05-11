package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/output"
	"github.com/yourusername/vaultdiff/internal/vault"
)

type config struct {
	pathA      string
	pathB      string
	format     string
	vaultAddr  string
	vaultToken string
}

func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.pathA, "a", "", "First Vault secret path (required)")
	flag.StringVar(&cfg.pathB, "b", "", "Second Vault secret path (required)")
	flag.StringVar(&cfg.format, "format", "text", "Output format: text or json")
	flag.StringVar(&cfg.vaultAddr, "addr", os.Getenv("VAULT_ADDR"), "Vault server address")
	flag.StringVar(&cfg.vaultToken, "token", os.Getenv("VAULT_TOKEN"), "Vault token")
	flag.Parse()
	return cfg
}

func run(cfg config) error {
	if cfg.pathA == "" || cfg.pathB == "" {
		return fmt.Errorf("both -a and -b paths are required")
	}

	fmt, err := output.ParseFormat(cfg.format)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}

	client, err := vault.NewClient(cfg.vaultAddr, cfg.vaultToken)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	secretA, err := client.ReadSecret(cfg.pathA)
	if err != nil {
		return fmt.Errorf("failed to read path %q: %w", cfg.pathA, err)
	}

	secretB, err := client.ReadSecret(cfg.pathB)
	if err != nil {
		return fmt.Errorf("failed to read path %q: %w", cfg.pathB, err)
	}

	changes := diff.Compare(secretA, secretB)
	report := output.NewReport(cfg.pathA, cfg.pathB, changes)

	if err := output.Write(fmt, report, os.Stdout); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

func main() {
	cfg := parseFlags()
	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
