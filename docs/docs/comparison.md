---
sidebar_position: 7
title: Comparison
description: How groundctl compares to chezmoi, mise, asdf, and manual setup.
---

# groundctl vs Alternatives

groundctl fills a gap that no existing tool covers. Here's how it compares.

## Quick Comparison

| Feature | groundctl | chezmoi | mise/asdf | Ansible | Manual |
|---------|-----------|---------|-----------|---------|--------|
| **Detects machine drift** | Yes | No | Partial | No | No |
| **Team-defined standard** | Yes | Per-user | Per-project | Yes | Wiki/README |
| **Auto-fix drift** | Yes | No | Yes (versions) | Yes | No |
| **Cross-platform** | Yes | Yes | Partial | Linux-focused | N/A |
| **Zero config start** | `ground init` | Complex | Per-tool `.tool-versions` | Playbook | Manual |
| **CI/CD integration** | Built-in | No | No | Overkill | Scripts |
| **Secret management** | Yes | Yes | No | Yes (vault) | .env files |
| **Single binary** | Yes | Yes | Plugin system | Python | N/A |
| **Setup time** | 2 minutes | 30+ minutes | 10+ minutes | Hours | Hours |

## Detailed Comparison

### vs chezmoi / dotbot / yadm

**What they do:** Manage individual dotfiles (`.bashrc`, `.vimrc`, etc.) across machines.

**What they don't do:** Detect whether your machine has the right _tools_ at the right _versions_. chezmoi manages config files, not installed software.

**groundctl is different:** It checks tool versions against a team standard and auto-fixes drift. It complements dotfile managers rather than replacing them.

### vs mise / asdf

**What they do:** Manage language runtime versions per-project (Node 20 for project A, Node 18 for project B).

**What they don't do:** Define a team-wide standard, detect drift in non-language tools (docker, terraform, git), or integrate with CI to enforce standards.

**groundctl is different:** It covers the full machine environment (not just language runtimes), defines team standards (not per-project versions), and provides CI enforcement. groundctl can work alongside mise/asdf.

### vs Ansible / Puppet / Chef

**What they do:** Full infrastructure automation for servers and cloud instances.

**What they don't do (well):** Manage local developer machines. They're designed for remote provisioning, not for checking "does Sarah have the right Node version?"

**groundctl is different:** It's purpose-built for local developer machines, runs in seconds (not minutes), requires no agent or server, and produces human-friendly output.

### vs Manual Setup (Wiki / README)

**What teams do today:** Write a "Dev Setup" wiki page or README section. New hires follow it manually. It gets outdated. Things break.

**groundctl is different:** The standard is code (`.ground.yaml`), it's versioned in git, drift is detected automatically, and fixes are one command away. It's the executable version of your setup wiki.

## When to Use groundctl

Use groundctl when:

- Your team has a standard set of tools and versions
- New hires spend time debugging environment differences
- "Works on my machine" is a recurring problem
- You want CI to enforce environment standards
- You need secret references without plaintext in git

Don't use groundctl when:

- You only need dotfile syncing (use chezmoi)
- You only need per-project language versions (use mise/asdf)
- You're provisioning cloud servers (use Ansible/Terraform)
