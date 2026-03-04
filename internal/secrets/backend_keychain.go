package secrets

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// KeychainBackend resolves secrets from the OS credential store.
// Reference format: ${keychain://service/account}
// On macOS: uses `security find-generic-password`
// On Linux: uses `secret-tool lookup`
// On Windows: uses `cmdkey` / PowerShell Get-StoredCredential
type KeychainBackend struct{}

func (b *KeychainBackend) Name() string { return "keychain" }

func (b *KeychainBackend) Resolve(path string) (string, error) {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("keychain: path must be service/account, got %q", path)
	}
	service := parts[0]
	account := parts[1]

	switch runtime.GOOS {
	case "darwin":
		return b.resolveMacOS(service, account)
	case "linux":
		return b.resolveLinux(service, account)
	case "windows":
		return b.resolveWindows(service, account)
	default:
		return "", fmt.Errorf("keychain: unsupported OS %q", runtime.GOOS)
	}
}

func (b *KeychainBackend) resolveMacOS(service, account string) (string, error) {
	out, err := exec.Command("security", "find-generic-password",
		"-s", service, "-a", account, "-w").Output()
	if err != nil {
		return "", fmt.Errorf("keychain: could not find %s/%s in macOS Keychain: %w", service, account, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (b *KeychainBackend) resolveLinux(service, account string) (string, error) {
	out, err := exec.Command("secret-tool", "lookup",
		"service", service, "account", account).Output()
	if err != nil {
		return "", fmt.Errorf("keychain: could not find %s/%s via secret-tool: %w", service, account, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (b *KeychainBackend) resolveWindows(service, account string) (string, error) {
	// Use PowerShell to read from Windows Credential Manager
	script := fmt.Sprintf(
		`$cred = Get-StoredCredential -Target '%s/%s'; if ($cred) { $cred.GetNetworkCredential().Password } else { exit 1 }`,
		service, account)
	out, err := exec.Command("powershell", "-NoProfile", "-Command", script).Output()
	if err != nil {
		return "", fmt.Errorf("keychain: could not find %s/%s in Windows Credential Manager: %w", service, account, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (b *KeychainBackend) Check() error {
	switch runtime.GOOS {
	case "darwin":
		if _, err := exec.LookPath("security"); err != nil {
			return fmt.Errorf("macOS Keychain: 'security' command not found")
		}
	case "linux":
		if _, err := exec.LookPath("secret-tool"); err != nil {
			return fmt.Errorf("libsecret: 'secret-tool' not found. Install: sudo apt install libsecret-tools")
		}
	case "windows":
		// PowerShell is always available on Windows
	default:
		return fmt.Errorf("keychain backend not supported on %s", runtime.GOOS)
	}
	return nil
}
