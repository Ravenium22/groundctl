---
sidebar_position: 3
title: CLI Reference
description: Complete reference for all groundctl commands and flags.
---

# CLI Reference

## Global Flags

These flags are available on all commands:

| Flag | Description |
|------|-------------|
| `-v, --verbose` | Verbose output |
| `--debug` | Debug mode (show detection commands and timing) |
| `-h, --help` | Help for any command |

## Commands

### `ground init`

Scan your machine and create a `.ground.yaml` file.

```bash
ground init [flags]
```

| Flag | Description |
|------|-------------|
| `--force` | Overwrite existing `.ground.yaml` |
| `--dir <path>` | Directory to create config in (default: current) |

### `ground check`

Compare your machine against the team standard.

```bash
ground check [flags]
```

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to `.ground.yaml` (default: auto-detect) |
| `--json` | Output report as JSON |
| `-q, --quiet` | Only show errors and warnings |
| `--ci` | Output GitHub Actions annotations |

**Exit codes:** `0` = all ok, `1` = warnings only, `2` = errors present.

### `ground fix`

Auto-fix detected drift using available package managers.

```bash
ground fix [flags]
```

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to `.ground.yaml` |
| `--dry-run` | Show what would be done without doing it |
| `--auto` | Skip confirmation prompts |
| `--json` | Output fix plan as JSON |

### `ground snapshot`

Capture your machine state as JSON.

```bash
ground snapshot [flags]
```

| Flag | Description |
|------|-------------|
| `-o, --output <path>` | Write snapshot to file |

### `ground diff`

Show what changed since the last snapshot.

```bash
ground diff [flags]
```

| Flag | Description |
|------|-------------|
| `-f, --from <path>` | Path to previous snapshot JSON |
| `--json` | Output diff as JSON |

### `ground report`

Generate a drift report in various formats.

```bash
ground report [flags]
```

| Flag | Description |
|------|-------------|
| `-f, --format <fmt>` | Output format: `json`, `markdown`, `html` |
| `-o, --output <path>` | Write to file (default: stdout) |
| `-c, --config <path>` | Path to `.ground.yaml` |

### `ground export`

Export environment spec for use in other contexts.

```bash
ground export [flags]
```

| Flag | Description |
|------|-------------|
| `-f, --format <fmt>` | Output format: `dockerfile`, `shell`, `json` |
| `-o, --output <path>` | Write to file (default: stdout) |
| `-c, --config <path>` | Path to `.ground.yaml` |

### `ground doctor`

Diagnose groundctl configuration and environment.

```bash
ground doctor
```

Checks: config validity, git availability, package managers, profiles, shell.

### `ground watch`

Watch for drift changes in the background.

```bash
ground watch [flags]
```

| Flag | Description |
|------|-------------|
| `-i, --interval <dur>` | Check interval (default: `30s`) |
| `-c, --config <path>` | Path to `.ground.yaml` |

Press `Ctrl+C` to stop.

### `ground pull`

Fetch a team standard from a remote source.

```bash
ground pull <source> [flags]
```

| Flag | Description |
|------|-------------|
| `-o, --output <path>` | Output file path |

Sources: git repos, URLs, local file paths.

### `ground push`

Publish your config to the team repo.

```bash
ground push [flags]
```

| Flag | Description |
|------|-------------|
| `-m, --message <msg>` | Commit message |

### `ground profile`

Manage environment profiles.

```bash
ground profile <subcommand>
```

| Subcommand | Description |
|------------|-------------|
| `list` | List all saved profiles |
| `save <name>` | Save current config as a profile |
| `switch <name>` | Set active profile |
| `show <name>` | Display a profile's config |
| `delete <name>` | Remove a saved profile |

### `ground secrets`

Manage secret references in your config.

```bash
ground secrets <subcommand> [flags]
```

| Subcommand | Description |
|------------|-------------|
| `check` | Validate all secret references are resolvable |
| `list` | List all configured secrets |
| `env` | Generate a `.env` file from secret references |

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to `.ground.yaml` |
| `-o, --output <path>` | Output path for `.env` (default: `.env`) |

### `ground validate`

Validate a `.ground.yaml` file.

```bash
ground validate [flags]
```

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to `.ground.yaml` |
| `--json` | Output errors as JSON |

### `ground hook`

Generate shell hooks for auto-checking on `cd`.

```bash
ground hook <shell>
ground hook --prompt <prompt>
```

Shells: `bash`, `zsh`, `fish`, `powershell`.
Prompts: `starship`, `p10k`.

### `ground completion`

Generate tab completion scripts.

```bash
ground completion <shell>
```

Shells: `bash`, `zsh`, `fish`, `powershell`.

### `ground version`

Print the groundctl version.

```bash
ground version
```
