# How groundctl Works Under the Hood

## Architecture of a Simple Tool

groundctl is deliberately simple. It's a single Go binary with no runtime dependencies. Here's how it works.

## The Detection Engine

At the core is a detection engine that knows how to find 18 common development tools. For each tool, it runs a version command (like `node --version`) and extracts the version number with a regex.

Detection runs in parallel using goroutines, bounded by CPU count. On a typical machine, all 18 tools are detected in under 500ms. Results are cached for 5 minutes to make repeated checks near-instant.

## The Diff Engine

The diff engine compares detected versions against constraints in `.ground.yaml`. It uses semver comparison (via the Masterminds/semver library) to evaluate constraints like `>=20.0.0`, `^3.11`, and `~1.6.0`.

Each tool gets a status: ok, warning, or error. The exit code reflects the worst status (0 = all ok, 1 = warnings, 2 = errors), making it CI-friendly.

## The Fix Engine

When drift is detected, the fix engine determines the best remediation. It detects which package managers are available (brew, apt, winget, scoop, choco, dnf, pacman) and picks the best one for each tool on the current OS.

Fix runs in dry-run mode by default, showing what it would do before doing it. The `--auto` flag skips confirmation for CI environments.

## Secret Management

Secrets use a `${backend://path}` reference syntax. Four backends are supported: environment variables, 1Password CLI, HashiCorp Vault, and OS keychain. The `ground secrets env` command resolves all references and generates a `.env` file.

Secret values are never displayed in terminal output — they're always masked (e.g., `sk*********45`).

## Caching

Detection results are cached in `~/.groundctl/cache/detect.json`. Cache entries have a 5-minute TTL and are invalidated when a tool's binary modification time changes. This means the first `ground check` takes ~500ms, but subsequent checks are near-instant.

## The Stack

- **Go** — single static binary, fast startup, cross-platform
- **Cobra** — CLI framework (powers kubectl, Hugo, GitHub CLI)
- **lipgloss** — terminal styling (by Charm)
- **Masterminds/semver** — version constraint evaluation
- **gopkg.in/yaml.v3** — YAML parsing

No CGo. No external dependencies at runtime. Works everywhere Go compiles.
