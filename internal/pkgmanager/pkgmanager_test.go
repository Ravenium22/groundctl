package pkgmanager

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetect(t *testing.T) {
	managers := Detect()
	// On any dev machine, at least one PM should be available
	// but we don't assert count since it depends on the machine
	for _, m := range managers {
		assert.NotEmpty(t, m.Name)
		assert.NotEmpty(t, m.Path)
		assert.Equal(t, runtime.GOOS, m.Platform)
	}
}

func TestListKnown(t *testing.T) {
	names := ListKnown()
	assert.NotEmpty(t, names, "should know at least one PM for current platform")
}

func TestListKnownForPlatform(t *testing.T) {
	tests := []struct {
		platform string
		expected []string
	}{
		{"darwin", []string{"brew"}},
		{"linux", []string{"apt", "dnf", "pacman", "brew"}},
		{"windows", []string{"winget", "scoop", "choco"}},
		{"freebsd", nil},
	}

	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			names := ListKnownForPlatform(tt.platform)
			if tt.expected == nil {
				assert.Empty(t, names)
			} else {
				assert.Equal(t, tt.expected, names)
			}
		})
	}
}

func TestDetectForPlatformUnknown(t *testing.T) {
	managers := DetectForPlatform("plan9")
	assert.Empty(t, managers)
}
