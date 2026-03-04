package secrets

import (
	"fmt"
	"os"
)

// EnvBackend resolves secrets from environment variables.
// Reference format: ${env://VAR_NAME}
type EnvBackend struct{}

func (b *EnvBackend) Name() string { return "env" }

func (b *EnvBackend) Resolve(path string) (string, error) {
	val, ok := os.LookupEnv(path)
	if !ok {
		return "", fmt.Errorf("environment variable %q not set", path)
	}
	return val, nil
}

func (b *EnvBackend) Check() error {
	return nil // env is always available
}
