---
sidebar_position: 4
title: Configuration
description: Full reference for .ground.yaml configuration.
---

# Configuration

groundctl uses `.ground.yaml` files to define environment standards. This page covers every option.

## File Location

groundctl searches for `.ground.yaml` starting from the current directory, walking up to the filesystem root. This lets you have per-project configs in repo roots.

## Full Schema

```yaml
# Optional metadata
name: my-project
description: Team development environment standard

# Inherit from a saved profile
extends: base-profile

# Team sharing metadata
team:
  org: acme-corp
  repo: github.com/acme/groundfile
  branch: main

# Required and recommended tools
tools:
  - name: node
    version: ">=20.0.0"
    severity: required
    install_cmd: "brew install node"

  - name: terraform
    version: "^1.6"
    severity: recommended

# Secret references (never stored in plaintext)
secrets:
  - name: DATABASE_URL
    ref: "${env://DATABASE_URL}"
    description: PostgreSQL connection string

  - name: API_KEY
    ref: "${op://Engineering/api-key/credential}"
    description: Production API key
```

## Fields

### Top-Level

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Config name (for display) |
| `description` | string | Human-readable description |
| `extends` | string | Parent profile name for inheritance |
| `team` | object | Team sharing metadata |
| `tools` | list | Tool requirements |
| `secrets` | list | Secret references |

### Tool Spec

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Tool identifier (e.g. `node`, `docker`) |
| `version` | string | No | Semver constraint (e.g. `>=20.0.0`, `^3.11`) |
| `severity` | string | No | `required` or `recommended` (default: `required`) |
| `install_cmd` | string | No | Custom install command hint |

### Secret Spec

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Environment variable name |
| `ref` | string | Yes | Secret reference (`${backend://path}`) |
| `description` | string | No | Human-readable description |

## Supported Tools

groundctl detects 18 tools out of the box:

`node`, `npm`, `python`, `pip`, `go`, `rustc`, `cargo`, `java`, `ruby`, `make`, `docker`, `docker-compose`, `kubectl`, `terraform`, `git`, `gh`, `curl`, `wget`

## Version Constraint Syntax

| Pattern | Meaning | Example |
|---------|---------|---------|
| `>=1.2.0` | Greater than or equal | `>=20.0.0` |
| `^1.2.0` | Compatible release (same major) | `^3.11` |
| `~1.2.0` | Approximately (same minor) | `~1.6.0` |
| `1.2.0` | Exact match | `3.12.1` |
| `*` | Any version | Accept anything |
| `>=1.0.0 <2.0.0` | Range | `>=1.6.0 <2.0.0` |

## Inheritance

When `extends` is set, tools from the parent profile are merged:

- Child tools override parent tools with the same name
- Parent tools not in the child are inherited
- Inheritance chains up to 10 levels deep are supported

## Validation

Run `ground validate` to check your config for issues:

```bash
ground validate
ground validate --json
```

Checks: empty tools, missing names, duplicate tools, invalid severity, invalid version constraints, invalid secret references, duplicate secrets.
