package secrets

import (
	"fmt"
	"os/exec"
	"strings"
)

// VaultBackend resolves secrets from HashiCorp Vault.
// Reference format: ${vault://secret/path#field}
// Uses the path before # as the secret path, and after # as the field name.
// If no # is present, returns the first field value.
// Requires the `vault` CLI to be installed and VAULT_ADDR set.
type VaultBackend struct{}

func (b *VaultBackend) Name() string { return "vault" }

func (b *VaultBackend) Resolve(path string) (string, error) {
	secretPath := path
	field := ""
	if idx := strings.Index(path, "#"); idx >= 0 {
		secretPath = path[:idx]
		field = path[idx+1:]
	}

	args := []string{"kv", "get", "-mount=secret"}
	if field != "" {
		args = append(args, fmt.Sprintf("-field=%s", field))
	}
	args = append(args, secretPath)

	out, err := exec.Command("vault", args...).Output()
	if err != nil {
		return "", fmt.Errorf("vault: could not read %q: %w", path, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (b *VaultBackend) Check() error {
	if _, err := exec.LookPath("vault"); err != nil {
		return fmt.Errorf("HashiCorp Vault CLI not found. Install: https://developer.hashicorp.com/vault/install")
	}
	if err := exec.Command("vault", "status").Run(); err != nil {
		return fmt.Errorf("vault not reachable or sealed. Check VAULT_ADDR and run: vault login")
	}
	return nil
}
