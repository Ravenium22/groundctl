package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/groundctl/groundctl/internal/config"
	"github.com/groundctl/groundctl/internal/model"
)

// Dir returns the profiles directory (~/.groundctl/profiles/).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".groundctl", "profiles"), nil
}

// EnsureDir creates the profiles directory if it doesn't exist.
func EnsureDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("could not create profiles directory: %w", err)
	}
	return dir, nil
}

// List returns all available profile names.
func List() ([]string, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not read profiles directory: %w", err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			names = append(names, strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml"))
		}
	}
	return names, nil
}

// path returns the file path for a named profile.
func path(name string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".yaml"), nil
}

// Load reads a named profile from the profiles directory.
func Load(name string) (*model.GroundConfig, error) {
	p, err := path(name)
	if err != nil {
		return nil, err
	}
	return config.Load(p)
}

// Save writes a config as a named profile.
func Save(name string, cfg *model.GroundConfig) error {
	dir, err := EnsureDir()
	if err != nil {
		return err
	}
	p := filepath.Join(dir, name+".yaml")
	return config.Save(p, cfg)
}

// Exists checks if a named profile exists.
func Exists(name string) bool {
	p, err := path(name)
	if err != nil {
		return false
	}
	return config.Exists(p)
}

// Delete removes a named profile.
func Delete(name string) error {
	p, err := path(name)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not delete profile %q: %w", name, err)
	}
	return nil
}

// Resolve loads a profile and resolves inheritance via the "extends" field.
// Returns the fully merged config.
func Resolve(name string) (*model.GroundConfig, error) {
	return resolveWithDepth(name, 0)
}

const maxInheritanceDepth = 10

func resolveWithDepth(name string, depth int) (*model.GroundConfig, error) {
	if depth > maxInheritanceDepth {
		return nil, fmt.Errorf("profile inheritance chain too deep (max %d): possible cycle", maxInheritanceDepth)
	}

	cfg, err := Load(name)
	if err != nil {
		return nil, fmt.Errorf("could not load profile %q: %w", name, err)
	}

	if cfg.Extends == "" {
		return cfg, nil
	}

	parent, err := resolveWithDepth(cfg.Extends, depth+1)
	if err != nil {
		return nil, fmt.Errorf("could not resolve parent profile %q: %w", cfg.Extends, err)
	}

	cfg.MergeTools(parent)
	return cfg, nil
}

// ActiveProfileFile returns the path to the file that stores the active profile name.
func ActiveProfileFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".groundctl", "active-profile"), nil
}

// GetActive returns the currently active profile name, or "" if none.
func GetActive() string {
	f, err := ActiveProfileFile()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(f)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// SetActive sets the active profile name.
func SetActive(name string) error {
	f, err := ActiveProfileFile()
	if err != nil {
		return err
	}
	dir := filepath.Dir(f)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(f, []byte(name+"\n"), 0644)
}

// ClearActive removes the active profile setting.
func ClearActive() error {
	f, err := ActiveProfileFile()
	if err != nil {
		return err
	}
	if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
