---
sidebar_position: 2
title: Team Setup
description: Set up groundctl for your engineering team in 10 minutes.
---

# Team Setup

Set up groundctl for your engineering team in 10 minutes. By the end, every developer can sync to the team standard with a single command.

## 1. Create the Team Standard

Start by initializing in your shared repo:

```bash
ground init
```

Edit `.ground.yaml` to define your team's requirements:

```yaml
name: acme-engineering
description: Acme Corp standard dev environment
tools:
  - name: node
    version: ">=20.0.0"
    severity: required
  - name: python
    version: ">=3.11"
    severity: required
  - name: docker
    version: ">=24.0.0"
    severity: required
  - name: terraform
    version: "^1.6"
    severity: recommended
  - name: gh
    version: ">=2.0.0"
    severity: recommended
```

### Severity Levels

- **`required`** - Missing or wrong version is an error (exit code 2)
- **`recommended`** - Missing or wrong version is a warning (exit code 1)

### Version Constraints

| Syntax | Meaning | Example |
|--------|---------|---------|
| `>=1.2.0` | At least this version | `>=20.0.0` |
| `^1.2.0` | Compatible (same major) | `^3.11` matches 3.11-3.x |
| `~1.2.0` | Patch-level changes only | `~1.6.0` matches 1.6.x |
| `*` | Any version | Accept anything installed |

## 2. Share with the Team

Commit `.ground.yaml` to your repo:

```bash
git add .ground.yaml
git commit -m "Add team environment standard"
git push
```

Team members pull the standard:

```bash
# From a git repo
ground pull github.com/acme/groundfile

# From a URL
ground pull https://raw.githubusercontent.com/acme/repo/main/.ground.yaml
```

## 3. Use Profiles

Developers can maintain multiple profiles for different projects:

```bash
# Save the current config as a profile
ground profile save work

# Switch between profiles
ground profile switch personal
ground profile switch work

# List all profiles
ground profile list
```

### Profile Inheritance

Create a personal config that extends the team standard:

```yaml
name: my-env
extends: work
tools:
  - name: neovim
    version: ">=0.9.0"
    severity: recommended
```

Child tools override parent tools with the same name; parent tools are inherited for tools not listed in the child.

## 4. Enforce in CI

Add groundctl to your CI pipeline to catch drift before it reaches production.

### GitHub Actions

```yaml
# .github/workflows/env-check.yml
name: Environment Check
on: [pull_request]

jobs:
  ground-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: Ravenium22/groundctl/.github/actions/setup-ground@main
        with:
          fail-on-drift: 'true'
          report-format: 'markdown'
```

### Manual CI Setup

```bash
# Install
curl -fsSL https://raw.githubusercontent.com/Ravenium22/groundctl/main/install.sh | sh

# Check with CI annotations
ground check --ci

# Generate a report
ground report --format markdown --output drift-report.md
```

## 5. Shell Integration

Auto-check when entering project directories:

```bash
# Bash - add to ~/.bashrc
eval "$(ground hook bash)"

# Zsh - add to ~/.zshrc
eval "$(ground hook zsh)"

# Fish - add to ~/.config/fish/config.fish
ground hook fish | source
```

Add a drift indicator to your prompt:

```bash
# Starship
ground hook --prompt starship >> ~/.config/starship.toml

# Powerlevel10k
ground hook --prompt p10k
```

## What's Next?

- **[Secrets](secrets)** - Manage API keys and tokens safely in your config
- **[CLI Reference](cli-reference)** - Full command documentation
- **[CI/CD Integration](ci-cd)** - Advanced pipeline setup
