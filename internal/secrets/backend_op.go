package secrets

import (
	"fmt"
	"os/exec"
	"strings"
)

// OpBackend resolves secrets from 1Password CLI.
// Reference format: ${op://vault/item/field}
// Requires the `op` CLI to be installed and authenticated.
type OpBackend struct{}

func (b *OpBackend) Name() string { return "op" }

func (b *OpBackend) Resolve(path string) (string, error) {
	ref := "op://" + path
	out, err := exec.Command("op", "read", ref).Output()
	if err != nil {
		return "", fmt.Errorf("1password: could not read %q: %w", ref, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (b *OpBackend) Check() error {
	if _, err := exec.LookPath("op"); err != nil {
		return fmt.Errorf("1Password CLI (op) not found. Install: https://1password.com/downloads/command-line/")
	}
	// Check if signed in by running `op whoami`
	if err := exec.Command("op", "whoami").Run(); err != nil {
		return fmt.Errorf("1Password CLI not authenticated. Run: op signin")
	}
	return nil
}
