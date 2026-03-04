---
sidebar_position: 5
title: Secrets
description: Manage API keys, tokens, and credentials safely in your config.
---

# Secrets

groundctl lets you reference secrets in `.ground.yaml` without storing them in plaintext. Secrets are resolved at runtime from external providers.

## Secret Reference Syntax

```
${backend://path}
```

| Backend | Syntax | Provider |
|---------|--------|----------|
| `env` | `${env://VAR_NAME}` | Environment variables |
| `op` | `${op://vault/item/field}` | 1Password CLI |
| `vault` | `${vault://secret/path#field}` | HashiCorp Vault |
| `keychain` | `${keychain://service/account}` | OS credential store |

## Configuration

Add secrets to your `.ground.yaml`:

```yaml
secrets:
  - name: DATABASE_URL
    ref: "${env://DATABASE_URL}"
    description: PostgreSQL connection string

  - name: API_KEY
    ref: "${op://Engineering/api-key/credential}"
    description: Production API key

  - name: VAULT_TOKEN
    ref: "${vault://secret/myapp#token}"
    description: Service token from Vault

  - name: SIGNING_KEY
    ref: "${keychain://myapp/signing-key}"
    description: Code signing key from OS keychain
```

## Commands

### Check secrets

Validate that all references can be resolved:

```bash
ground secrets check
```

### List secrets

Show all configured secret references:

```bash
ground secrets list
```

### Generate .env file

Resolve all secrets and write to a `.env` file:

```bash
ground secrets env
ground secrets env --output .env.local
```

The `.env` file is written with `0600` permissions. Secret values are masked in terminal output.

## Backend Setup

### Environment Variables (`env`)

No setup required. References environment variables on the current system.

### 1Password (`op`)

Requires the [1Password CLI](https://1password.com/downloads/command-line/):

```bash
# Install
brew install 1password-cli

# Authenticate
op signin
```

### HashiCorp Vault (`vault`)

Requires the [Vault CLI](https://developer.hashicorp.com/vault/install):

```bash
# Install
brew install vault

# Configure
export VAULT_ADDR=https://vault.example.com
vault login
```

Use `#field` to select a specific field: `${vault://secret/db#password}`

### OS Keychain (`keychain`)

Uses the platform credential store:

- **macOS**: Keychain Access (`security` command)
- **Linux**: libsecret (`secret-tool` command)
- **Windows**: Credential Manager (PowerShell)

## Security Model

- Secret values are **never** written to `.ground.yaml` or any config file
- Terminal output always shows **masked** values (e.g. `sk*********45`)
- `.env` files are created with restrictive permissions (`0600`)
- Add `.env` to your `.gitignore` to prevent accidental commits
