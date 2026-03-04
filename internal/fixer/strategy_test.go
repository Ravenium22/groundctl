package fixer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveStrategy_BrewInstall(t *testing.T) {
	plan := ResolveStrategyForPlatform("node", ActionInstall, []string{"brew"}, ">=18.0.0", "", "darwin")
	assert.Equal(t, "brew", plan.Manager)
	assert.Equal(t, ActionInstall, plan.Action)
	assert.Equal(t, []string{"brew", "install", "node"}, plan.Command)
	assert.Equal(t, "brew install node", plan.CommandStr)
	assert.Empty(t, plan.ManualHint)
}

func TestResolveStrategy_AptInstall(t *testing.T) {
	plan := ResolveStrategyForPlatform("python", ActionInstall, []string{"apt"}, ">=3.10", "", "linux")
	assert.Equal(t, "apt", plan.Manager)
	assert.Equal(t, []string{"sudo", "apt-get", "install", "-y", "python3"}, plan.Command)
}

func TestResolveStrategy_WingetInstall(t *testing.T) {
	plan := ResolveStrategyForPlatform("go", ActionInstall, []string{"winget"}, ">=1.21", "", "windows")
	assert.Equal(t, "winget", plan.Manager)
	assert.Contains(t, plan.Command, "GoLang.Go")
}

func TestResolveStrategy_Upgrade(t *testing.T) {
	plan := ResolveStrategyForPlatform("node", ActionUpgrade, []string{"brew"}, ">=20.0.0", "18.19.0", "darwin")
	assert.Equal(t, ActionUpgrade, plan.Action)
	assert.Equal(t, []string{"brew", "upgrade", "node"}, plan.Command)
	assert.Equal(t, "18.19.0", plan.Current)
}

func TestResolveStrategy_PriorityOrder(t *testing.T) {
	// winget should be preferred over scoop on Windows
	plan := ResolveStrategyForPlatform("git", ActionInstall, []string{"winget", "scoop", "choco"}, "", "", "windows")
	assert.Equal(t, "winget", plan.Manager)
}

func TestResolveStrategy_FallbackToPM(t *testing.T) {
	// First PM doesn't have the tool, second does
	plan := ResolveStrategyForPlatform("node", ActionInstall, []string{"dnf"}, ">=18.0.0", "", "linux")
	assert.Equal(t, "dnf", plan.Manager)
	assert.Contains(t, plan.Command, "nodejs")
}

func TestResolveStrategy_NoManagerAvailable(t *testing.T) {
	plan := ResolveStrategyForPlatform("node", ActionInstall, []string{}, ">=18.0.0", "", "linux")
	assert.Empty(t, plan.Manager)
	assert.Empty(t, plan.Command)
	assert.NotEmpty(t, plan.ManualHint)
}

func TestResolveStrategy_UnknownTool(t *testing.T) {
	plan := ResolveStrategyForPlatform("unkown-tool-xyz", ActionInstall, []string{"brew"}, "", "", "darwin")
	assert.Empty(t, plan.Manager)
	assert.Contains(t, plan.ManualHint, "No install strategy known")
}

func TestResolveStrategy_ManualHintForKnownTool(t *testing.T) {
	// docker with no PMs should give a useful hint
	plan := ResolveStrategyForPlatform("docker", ActionInstall, []string{}, ">=24.0.0", "", "linux")
	assert.Contains(t, plan.ManualHint, "docker")
}

func TestResolveStrategy_AllPMs(t *testing.T) {
	pms := []string{"brew", "apt", "winget", "scoop", "choco", "dnf", "pacman"}
	for _, pm := range pms {
		plan := ResolveStrategyForPlatform("git", ActionInstall, []string{pm}, "", "", "linux")
		assert.Equal(t, pm, plan.Manager, "should resolve for PM: %s", pm)
		assert.NotEmpty(t, plan.Command)
	}
}

func TestLookupPackageName(t *testing.T) {
	pkg, ok := LookupPackageName("node", "brew")
	assert.True(t, ok)
	assert.Equal(t, "node", pkg)

	pkg, ok = LookupPackageName("node", "apt")
	assert.True(t, ok)
	assert.Equal(t, "nodejs", pkg)

	_, ok = LookupPackageName("nonexistent", "brew")
	assert.False(t, ok)

	_, ok = LookupPackageName("node", "nonexistent-pm")
	assert.False(t, ok)
}

func TestFormatCommand(t *testing.T) {
	assert.Equal(t, "brew install node", formatCommand([]string{"brew", "install", "node"}))
	assert.Equal(t, "sudo apt-get install -y python3", formatCommand([]string{"sudo", "apt-get", "install", "-y", "python3"}))
}
