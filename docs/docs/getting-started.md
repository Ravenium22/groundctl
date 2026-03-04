---
sidebar_position: 1
title: Getting Started
description: Go from zero to your first drift check in under 5 minutes.
---

# Getting Started

Get from zero to your first `ground check` in under 5 minutes.

## Install

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<Tabs>
<TabItem value="curl" label="curl (macOS/Linux)" default>

```bash
curl -fsSL https://raw.githubusercontent.com/Ravenium22/groundctl/main/install.sh | sh
```

</TabItem>
<TabItem value="go" label="Go">

```bash
go install github.com/Ravenium22/groundctl@latest
```

</TabItem>
<TabItem value="docker" label="Docker">

```bash
docker run --rm ghcr.io/Ravenium22/groundctl version
```

</TabItem>
<TabItem value="manual" label="Manual Binary">

```bash
# Download the latest release archive for your OS/arch:
# https://github.com/Ravenium22/groundctl/releases/latest
```

</TabItem>
</Tabs>

Homebrew and winget packages are planned but not yet published from this repo.

Verify the installation:

```bash
ground version
```

## Step 1: Initialize

Run `ground init` in your project directory. groundctl scans your machine and creates a `.ground.yaml` file:

```bash
cd my-project
ground init
```

```
Scanning machine for installed tools...

Created .ground.yaml with 11 tools detected
  node    >=22.10.0
  npm     >=10.9.0
  python  >=3.13.7
  go      >=1.26.0
  git     >=2.47.1
  docker  >=24.0.7
  ...
```

## Step 2: Check for Drift

Run `ground check` to compare your machine against the standard:

```bash
ground check
```

```
groundctl drift report
config: .ground.yaml

  [ok]  node          20.10.0
  [ok]  python        3.12.1
  [ERR] docker        not found
  [!!]  terraform     version drift (have: 1.5.0, want: >=1.6.0)

--------------------------------------------------
  4 checked  2 ok  1 warnings  1 errors

  Run 'ground fix' to resolve drift.
```

## Step 3: Fix Drift

Run `ground fix` to automatically resolve detected issues:

```bash
# Preview what will happen
ground fix --dry-run

# Fix interactively (confirm each step)
ground fix

# Fix everything automatically
ground fix --auto
```

## Step 4: Commit and Share

Commit `.ground.yaml` to your repo so the whole team benefits:

```bash
git add .ground.yaml
git commit -m "Add groundctl environment standard"
git push
```

Your teammates can now run `ground check` to see how their machine compares.

## What's Next?

- **[Team Setup](team-setup)** - Set up profiles, sharing, and CI enforcement
- **[CLI Reference](cli-reference)** - Full command reference
- **[Configuration](configuration)** - Customize `.ground.yaml`
- **[Secrets](secrets)** - Manage secret references safely
