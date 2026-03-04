---
sidebar_position: 6
title: CI/CD Integration
description: Enforce environment standards in your CI/CD pipelines.
---

# CI/CD Integration

groundctl integrates with CI/CD pipelines to enforce environment standards and catch drift before it reaches production.

## GitHub Actions

### Official Action

Use the official `setup-ground` action:

```yaml
name: Environment Check
on: [pull_request]

jobs:
  ground-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: Ravenium22/groundctl/.github/actions/setup-ground@main
        with:
          version: 'latest'
          fail-on-drift: 'true'
          report-format: 'markdown'
```

#### Inputs

| Input | Default | Description |
|-------|---------|-------------|
| `version` | `latest` | groundctl version to install |
| `config-path` | auto-detect | Path to `.ground.yaml` |
| `fail-on-drift` | `true` | Fail the step if drift is detected |
| `report-format` | (empty) | Generate report: `json`, `markdown`, `html` |

### Manual Setup

```yaml
- name: Install groundctl
  run: curl -fsSL https://raw.githubusercontent.com/Ravenium22/groundctl/main/install.sh | sh

- name: Check environment
  run: ground check --ci
```

The `--ci` flag outputs GitHub Actions annotations that appear inline on PR diffs.

## GitLab CI

```yaml
ground-check:
  image: golang:latest
  script:
    - curl -fsSL https://raw.githubusercontent.com/Ravenium22/groundctl/main/install.sh | sh
    - ground check --json > drift-report.json
  artifacts:
    reports:
      dotenv: drift-report.json
    when: always
```

## CircleCI

```yaml
jobs:
  ground-check:
    docker:
      - image: cimg/go:1.25
    steps:
      - checkout
      - run:
          name: Install groundctl
          command: curl -fsSL https://raw.githubusercontent.com/Ravenium22/groundctl/main/install.sh | sh
      - run:
          name: Check environment
          command: ground check
```

## Reports

Generate drift reports in CI for review:

```bash
# Markdown (great for PR comments)
ground report --format markdown --output drift-report.md

# HTML (for artifact downloads)
ground report --format html --output drift-report.html

# JSON (for programmatic use)
ground report --format json --output drift-report.json
```

## Dockerfile Export

Generate Dockerfile stanzas from your config to keep CI images in sync:

```bash
ground export --format dockerfile --output Dockerfile.ground
```

This generates `ENV` and `RUN` stanzas for each tool in your `.ground.yaml`.

## Exit Codes

Use exit codes in CI scripts:

| Code | Meaning | CI Action |
|------|---------|-----------|
| `0` | All tools match | Pass |
| `1` | Warnings only | Pass or warn |
| `2` | Errors present | Fail |
