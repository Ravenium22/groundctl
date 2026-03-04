# groundctl v2 Roadmap

Published post-launch. This roadmap is shaped by beta feedback, launch-day conversations, and community requests.

## v1.1 — Quick Wins (2 weeks post-launch)

### Plugin System
- Custom tool detector plugins (Go or shell-based)
- `ground plugin install <name>` from a community registry
- Plugin template: `ground plugin init`
- Enables community-contributed detectors without core changes

### Config Inheritance Improvements
- Multiple `extends` (merge multiple team standards)
- Override precedence documentation
- `ground config merge` — combine two configs

### Quality of Life
- `ground upgrade` — self-update command
- `ground status` — one-line summary for scripts
- Colored diff in `ground diff` output
- `--format table` output option

## v1.2 — Enterprise Features (4 weeks post-launch)

### Policy Engine
- `ground policy` — enforce organizational constraints
- Blocklist: prevent specific tool versions (security CVEs)
- Allowlist: restrict to approved tools only
- Policy-as-code in `.ground-policy.yaml`

### Audit & Compliance
- `ground audit` — generate compliance report
- Machine fingerprinting for fleet tracking
- Export to CSV/PDF for compliance teams
- Integration with SIEM/logging platforms

### SSO & Team Management
- Team authentication via GitHub/GitLab org membership
- Role-based config access (admin, member, viewer)
- Config approval workflow (PR-based changes)

## v1.3 — Ecosystem Expansion (8 weeks post-launch)

### Language-Specific Modules
- Node.js: detect global packages, nvm/fnm config
- Python: detect virtualenv, pyenv, conda environments
- Ruby: detect rbenv/rvm, global gems
- Java: detect JAVA_HOME, Maven/Gradle config
- Rust: detect rustup toolchain, cargo global installs

### Container & Cloud Detection
- Docker Compose service versions
- Kubernetes context/cluster detection
- Cloud CLI tools (aws, gcloud, az) with profile detection
- Terraform workspace and provider version checking

### IDE Integration
- VS Code extension: inline drift indicators
- JetBrains plugin: drift panel
- Neovim plugin: statusline integration
- `ground export --format vscode-settings`

## v2.0 — The Platform (12 weeks post-launch)

### groundctl Cloud (Optional SaaS)
- Team dashboard: real-time fleet drift visibility
- Drift alerts via Slack/Teams/email
- Historical drift tracking and trends
- Config change audit log

### Auto-Remediation Engine
- Scheduled drift checks (cron-based)
- Auto-fix on detection (opt-in, per-tool)
- Rollback on failed auto-fix
- Fix dry-run notifications before applying

### Community Hub
- Shared config registry (like Homebrew taps)
- `ground search <tool>` — find community detectors
- Upvote/review system for shared configs
- "Teams using this config" social proof

## Principles

These guide all v2 decisions:

1. **CLI-first** — the CLI is always the primary interface; SaaS is optional
2. **Zero lock-in** — configs are YAML, portable, version-controlled
3. **Privacy by default** — telemetry opt-in, no data leaves the machine unless chosen
4. **Community-driven** — plugin system enables anyone to extend groundctl
5. **Cross-platform** — every feature works on macOS, Linux, and Windows
