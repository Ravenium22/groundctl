# Case Study: From 2-Day Onboarding to 10 Minutes

## The Situation

A 30-person engineering team running a microservices architecture. Tech stack: Node.js, Python, Go, Docker, Terraform, kubectl. New hires took an average of 2 days to get a working local environment.

## The Problem

- The setup wiki was 47 steps long and last updated 4 months ago
- 3 different Node versions across the team (18, 20, 22)
- Half the team had Docker Desktop, half had Colima
- Terraform version mismatch caused state file conflicts weekly
- Every sprint had at least one "works on my machine" blocker

## The Solution

The team adopted groundctl in three steps:

**Step 1: Baseline** (10 minutes)
The tech lead ran `ground init` and edited the resulting `.ground.yaml` to set team standards:

```yaml
name: platform-team
tools:
  - name: node
    version: ">=20.0.0"
    severity: required
  - name: docker
    version: ">=24.0.0"
    severity: required
  - name: terraform
    version: "^1.6"
    severity: required
  - name: kubectl
    version: ">=1.28.0"
    severity: required
```

**Step 2: Check** (2 minutes per developer)
Each developer ran `ground check`. The team found: 4 people had the wrong Node version, 2 were missing kubectl entirely, and 6 had outdated Terraform.

**Step 3: Fix** (30 seconds per developer)
`ground fix --auto` resolved every issue. CI was updated with `ground check --ci` to catch future drift.

## The Results

- **Onboarding time:** 2 days down to 10 minutes
- **"Works on my machine" incidents:** Dropped from ~4/sprint to 0
- **Terraform state conflicts:** Eliminated
- **Developer satisfaction:** "Why didn't we do this years ago?"

## The Workflow Now

1. `.ground.yaml` is committed to every repo
2. New hires run `ground check` then `ground fix` on day one
3. CI rejects PRs from machines that don't meet the standard
4. Shell hooks auto-check when developers `cd` into projects
