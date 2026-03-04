package pkgmanager

import (
	"os/exec"
	"runtime"
)

// Manager represents a package manager available on the system.
type Manager struct {
	Name     string // e.g. "brew", "apt", "winget"
	Path     string // resolved binary path
	Platform string // "darwin", "linux", "windows"
}

// registry of all known package managers keyed by platform.
var registry = map[string][]pmDef{
	"darwin": {
		{Name: "brew", Bin: "brew"},
	},
	"linux": {
		{Name: "apt", Bin: "apt-get"},
		{Name: "dnf", Bin: "dnf"},
		{Name: "pacman", Bin: "pacman"},
		{Name: "brew", Bin: "brew"},
	},
	"windows": {
		{Name: "winget", Bin: "winget"},
		{Name: "scoop", Bin: "scoop"},
		{Name: "choco", Bin: "choco"},
	},
}

type pmDef struct {
	Name string
	Bin  string
}

// Detect returns all package managers found on the current system,
// in priority order (first = preferred).
func Detect() []Manager {
	return DetectForPlatform(runtime.GOOS)
}

// DetectForPlatform returns package managers for a specific platform.
// Useful for testing.
func DetectForPlatform(platform string) []Manager {
	defs, ok := registry[platform]
	if !ok {
		return nil
	}

	var found []Manager
	for _, d := range defs {
		path, err := exec.LookPath(d.Bin)
		if err == nil {
			found = append(found, Manager{
				Name:     d.Name,
				Path:     path,
				Platform: platform,
			})
		}
	}
	return found
}

// ListKnown returns the names of all package managers we know about
// for the current platform.
func ListKnown() []string {
	return ListKnownForPlatform(runtime.GOOS)
}

// ListKnownForPlatform returns PM names for a specific platform.
func ListKnownForPlatform(platform string) []string {
	defs := registry[platform]
	names := make([]string, len(defs))
	for i, d := range defs {
		names[i] = d.Name
	}
	return names
}

// HasManager checks if a specific package manager is available.
func HasManager(name string) bool {
	for _, m := range Detect() {
		if m.Name == name {
			return true
		}
	}
	return false
}
