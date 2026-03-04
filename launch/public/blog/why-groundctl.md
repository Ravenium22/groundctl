# Why We Built groundctl

## The $50,000 Problem Every Engineering Team Ignores

Every engineering team has a dirty secret: their developers are all running slightly different environments. And it's costing them far more than they think.

## The Setup Tax

Picture this: a new engineer joins your team. They follow the setup wiki (last updated 8 months ago). Three hours in, they hit their first "command not found." Six hours in, they're debugging a version mismatch that nobody else has. Two days later, they finally have a working environment — but it's subtly different from everyone else's.

This isn't a new-hire problem. It's an ongoing tax. Someone upgrades Node. Someone installs a tool globally. Someone's on a different Terraform version. These differences accumulate silently until they surface as mysterious CI failures and "it works on my machine."

## The Gap in the Toolchain

We have great tools for managing individual pieces:
- **chezmoi/dotbot** manage dotfiles (config files, not tool versions)
- **mise/asdf** manage language runtimes per-project (not team-wide standards)
- **Ansible/Puppet** provision servers (overkill for local machines)

None of them answer the fundamental question: **"Does this machine match what the team expects?"**

## Enter groundctl

groundctl is built around a simple insight: environment drift is a diff problem. You have a standard (`.ground.yaml`), you have reality (your machine), and you need to see the difference.

```yaml
# .ground.yaml — committed to your repo
name: acme-backend
tools:
  - name: node
    version: ">=20.0.0"
    severity: required
  - name: docker
    version: ">=24.0.0"
    severity: required
  - name: terraform
    version: "^1.6"
    severity: recommended
```

```
$ ground check

  [ok]  node        22.10.0
  [ERR] docker      not found
  [!!]  terraform   version drift (have: 1.5.0, want: ^1.6)
```

That's the entire value proposition. One file defines the standard. One command shows the diff. One command fixes it.

## What's Next

groundctl is open source (MIT), written in Go, and works on macOS, Linux, and Windows. We'd love for you to try it and tell us what's missing.

`go install github.com/groundctl/groundctl@latest`
