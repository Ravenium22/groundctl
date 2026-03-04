package config

import (
	"fmt"
	"os"

	"github.com/Ravenium22/groundctl/internal/model"
	"gopkg.in/yaml.v3"
)

const DefaultConfigFile = ".ground.yaml"

// Load reads and parses a .ground.yaml file.
func Load(path string) (*model.GroundConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	var cfg model.GroundConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}
	return &cfg, nil
}

// Save writes a GroundConfig to a YAML file.
func Save(path string, cfg *model.GroundConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	header := []byte("# groundctl configuration\n# https://github.com/Ravenium22/groundctl\n\n")
	return os.WriteFile(path, append(header, data...), 0644)
}

// Exists checks if a config file exists at the given path.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
