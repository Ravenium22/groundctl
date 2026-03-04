package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashHook(t *testing.T) {
	hook := BashHook()
	assert.Contains(t, hook, "cd")
	assert.Contains(t, hook, ".ground.yaml")
	assert.Contains(t, hook, "ground check")
	assert.Contains(t, hook, "bashrc")
}

func TestZshHook(t *testing.T) {
	hook := ZshHook()
	assert.Contains(t, hook, "chpwd")
	assert.Contains(t, hook, ".ground.yaml")
	assert.Contains(t, hook, "ground check")
	assert.Contains(t, hook, "zshrc")
}

func TestFishHook(t *testing.T) {
	hook := FishHook()
	assert.Contains(t, hook, "PWD")
	assert.Contains(t, hook, ".ground.yaml")
	assert.Contains(t, hook, "ground check")
}

func TestPowerShellHook(t *testing.T) {
	hook := PowerShellHook()
	assert.Contains(t, hook, "Set-Location")
	assert.Contains(t, hook, ".ground.yaml")
	assert.Contains(t, hook, "ground check")
}

func TestStarshipSegment(t *testing.T) {
	seg := StarshipSegment()
	assert.Contains(t, seg, "starship.toml")
	assert.Contains(t, seg, "ground check")
	assert.Contains(t, seg, "custom.groundctl")
}

func TestP10kSegment(t *testing.T) {
	seg := P10kSegment()
	assert.Contains(t, seg, "p10k")
	assert.Contains(t, seg, "ground check")
}
