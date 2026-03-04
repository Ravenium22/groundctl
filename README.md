# groundctl

**`terraform plan` for your local developer machine.**

A CLI that tells your team exactly how their machine has drifted from the standard and fixes it with one command.

[![CI](https://github.com/Ravenium22/groundctl/actions/workflows/ci.yml/badge.svg)](https://github.com/Ravenium22/groundctl/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Ravenium22/groundctl)](https://goreportcard.com/report/github.com/Ravenium22/groundctl)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## The Problem

Every engineering team faces local environment drift. Someone updates their Node.js version, someone installs a CLI tool globally, and a new hire spends days debugging setup differences.

## Quick Start

### Install

```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/Ravenium22/groundctl/main/install.sh | sh

# Go (all platforms)
go install github.com/Ravenium22/groundctl@latest

# Docker
docker run --rm ghcr.io/ravenium22/groundctl version
```

`Homebrew` and `winget` packages are planned but not yet published from this repo.

### 2 Minutes to Your First Check

```bash
# 1. Initialize - scan your machine
ground init

# 2. Check - see what drifted
ground check

# 3. Fix - resolve drift automatically
ground fix
```

### Example Output

```
$ ground check

groundctl drift report
config: .ground.yaml

  [ok]  node          22.10.0
  [ok]  python        3.12.1
  [ERR] docker        not found
  [!!]  terraform     version drift (have: 1.5.0, want: >=1.6.0)

--------------------------------------------------
  4 checked  2 ok  1 warnings  1 errors

  Run 'ground fix' to resolve drift.
```

## Features

| Feature | Command |
|---------|---------|
| Detect drift | `ground check` |
| Auto-fix drift | `ground fix --auto` |
| Team profiles | `ground pull github.com/org/groundfile` |
| Secret management | `ground secrets env` |
| CI enforcement | `ground check --ci` |
| Shell hooks | `eval "$(ground hook bash)"` |
| Drift reports | `ground report --format html` |
| Self-diagnostic | `ground doctor` |

## `.ground.yaml`

```yaml
name: my-project
description: Team development environment standard
tools:
  - name: node
    version: ">=20.0.0"
    severity: required
  - name: python
    version: ">=3.11"
    severity: recommended
  - name: docker
    version: ">=24.0.0"
    severity: required
secrets:
  - name: API_KEY
    ref: "${op://Engineering/api-key/credential}"
    description: API key from 1Password
```

## Detected Tools

groundctl detects **18 tools** out of the box:

| Category | Tools |
|----------|-------|
| JavaScript | node, npm |
| Python | python, pip |
| Systems | go, rustc, cargo, java, ruby, make |
| DevOps | docker, docker-compose, kubectl, terraform |
| Utilities | git, gh, curl, wget |

## CI/CD Integration

### GitHub Actions

```yaml
- uses: Ravenium22/groundctl/.github/actions/setup-ground@main
  with:
    fail-on-drift: 'true'
    report-format: 'markdown'
```

### Export for Docker

```bash
ground export --format dockerfile
```

## Documentation

- [Getting Started](docs/docs/getting-started.md) - 5 minutes to first check
- [Team Setup](docs/docs/team-setup.md) - 10 minutes to team sync
- [CLI Reference](docs/docs/cli-reference.md) - all commands and flags
- [Configuration](docs/docs/configuration.md) - `.ground.yaml` reference
- [Secrets](docs/docs/secrets.md) - manage credentials safely
- [CI/CD Integration](docs/docs/ci-cd.md) - pipeline setup
- [Comparison](docs/docs/comparison.md) - groundctl vs alternatives

## Roadmap

### v1.0 (current)
- [x] Tool detection engine (18 tools, parallel detection)
- [x] Core commands: `ground init`, `ground snapshot`, `ground check`, `ground fix`
- [x] Team profiles and sharing (`ground pull/push`, `ground profile`)
- [x] Shell integration (`ground hook`, `ground completion`, `ground watch`, `ground diff`, `ground doctor`)
- [x] Secret-aware config (1Password, Vault, keychain, env backends)
- [x] CI/CD integration (`ground check --ci`, `ground report`, `ground export`, GitHub Action)
- [x] Cross-platform: macOS, Linux, Windows
- [x] Performance: parallel detection, result caching, <200ms checks
- [x] Opt-in telemetry, crash recovery, interactive demo

### v2.0 (planned)
- [ ] Plugin system for custom tool detectors
- [ ] Policy engine (blocklist/allowlist)
- [ ] Language-specific modules (node globals, python venvs, rust toolchains)
- [ ] IDE extensions (VS Code, JetBrains, Neovim)
- [ ] Optional SaaS dashboard for fleet visibility

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT
