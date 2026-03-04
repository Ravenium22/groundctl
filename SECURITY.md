# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 1.0.x   | Yes       |
| < 1.0   | No        |

## Reporting a Vulnerability

If you discover a security vulnerability in groundctl, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, email **security@groundctl.dev** with:

1. A description of the vulnerability
2. Steps to reproduce
3. Affected version(s)
4. Any potential impact assessment

We will acknowledge your report within 48 hours and provide a timeline for a fix within 5 business days.

## Security Model

### What groundctl accesses

- **Tool binaries**: Runs `--version` commands on locally installed tools (read-only)
- **Config files**: Reads/writes `.ground.yaml` and `~/.groundctl/` directory
- **Package managers**: Invokes package managers only during `ground fix` (with user confirmation)
- **Secret backends**: Resolves secret references via CLI tools (op, vault) — never stores secrets

### What groundctl does NOT do

- No network requests (except `ground pull` from user-specified git repos)
- No data collection without opt-in (`ground telemetry on`)
- No credential storage — secrets are resolved at runtime, never persisted
- No elevated privilege operations — `ground fix` uses the user's existing package manager permissions

### Telemetry

Telemetry is strictly opt-in. When enabled, only anonymous usage data is collected:
- Command name, OS, architecture
- Execution duration and exit code
- groundctl version

No personally identifiable information, tool versions, config contents, or secret references are ever collected.

## Dependencies

groundctl uses well-maintained, widely-adopted Go dependencies:
- `github.com/spf13/cobra` — CLI framework
- `github.com/Masterminds/semver/v3` — Version constraint parsing
- `github.com/charmbracelet/lipgloss` — Terminal styling
- `gopkg.in/yaml.v3` — YAML parsing

Dependencies are kept minimal and regularly audited via `go mod tidy` and Dependabot.
