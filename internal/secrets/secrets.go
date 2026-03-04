package secrets

import (
	"fmt"
	"regexp"
	"strings"
)

// SecretRef represents a parsed secret reference like ${op://vault/item/field}.
type SecretRef struct {
	Raw     string // original string, e.g. "${op://vault/item/field}"
	Backend string // backend name, e.g. "op", "vault", "keychain", "env"
	Path    string // backend-specific path, e.g. "vault/item/field"
}

// refPattern matches ${backend://path} references.
var refPattern = regexp.MustCompile(`\$\{(\w+)://([^}]+)\}`)

// ParseRef parses a single secret reference string.
// Returns nil if the string is not a valid secret reference.
func ParseRef(s string) *SecretRef {
	m := refPattern.FindStringSubmatch(s)
	if m == nil {
		return nil
	}
	return &SecretRef{
		Raw:     m[0],
		Backend: m[1],
		Path:    m[2],
	}
}

// ParseRefs extracts all secret references from a string.
// A string may contain multiple refs, e.g. "host=${vault://db/host}:${vault://db/port}".
func ParseRefs(s string) []SecretRef {
	matches := refPattern.FindAllStringSubmatch(s, -1)
	refs := make([]SecretRef, 0, len(matches))
	for _, m := range matches {
		refs = append(refs, SecretRef{
			Raw:     m[0],
			Backend: m[1],
			Path:    m[2],
		})
	}
	return refs
}

// Backend is the interface that secret provider backends must implement.
type Backend interface {
	// Name returns the backend identifier (e.g. "op", "vault", "keychain", "env").
	Name() string

	// Resolve retrieves the secret value for the given path.
	Resolve(path string) (string, error)

	// Check verifies the backend is available (CLI installed, authenticated, etc).
	Check() error
}

// Registry holds registered secret backends.
type Registry struct {
	backends map[string]Backend
}

// NewRegistry creates a new empty backend registry.
func NewRegistry() *Registry {
	return &Registry{
		backends: make(map[string]Backend),
	}
}

// Register adds a backend to the registry.
func (r *Registry) Register(b Backend) {
	r.backends[b.Name()] = b
}

// Get returns the backend for the given name, or an error if not found.
func (r *Registry) Get(name string) (Backend, error) {
	b, ok := r.backends[name]
	if !ok {
		return nil, fmt.Errorf("unknown secret backend %q", name)
	}
	return b, nil
}

// Resolve resolves a single secret reference to its value.
func (r *Registry) Resolve(ref SecretRef) (string, error) {
	b, err := r.Get(ref.Backend)
	if err != nil {
		return "", err
	}
	return b.Resolve(ref.Path)
}

// ResolveString replaces all secret references in a string with their values.
func (r *Registry) ResolveString(s string) (string, error) {
	refs := ParseRefs(s)
	if len(refs) == 0 {
		return s, nil
	}

	result := s
	for _, ref := range refs {
		val, err := r.Resolve(ref)
		if err != nil {
			return "", fmt.Errorf("resolving %s: %w", ref.Raw, err)
		}
		result = strings.Replace(result, ref.Raw, val, 1)
	}
	return result, nil
}

// CheckRef verifies that a single secret reference can be resolved.
func (r *Registry) CheckRef(ref SecretRef) error {
	b, err := r.Get(ref.Backend)
	if err != nil {
		return err
	}
	if err := b.Check(); err != nil {
		return fmt.Errorf("backend %q not available: %w", ref.Backend, err)
	}
	_, err = b.Resolve(ref.Path)
	return err
}

// DefaultRegistry returns a registry with all built-in backends registered.
func DefaultRegistry() *Registry {
	r := NewRegistry()
	r.Register(&EnvBackend{})
	r.Register(&OpBackend{})
	r.Register(&VaultBackend{})
	r.Register(&KeychainBackend{})
	return r
}

// MaskValue replaces a secret value with a masked version for display.
func MaskValue(val string) string {
	if len(val) <= 4 {
		return "****"
	}
	return val[:2] + strings.Repeat("*", len(val)-4) + val[len(val)-2:]
}

// MaskString replaces all occurrences of known secret values in a string.
func MaskString(s string, secrets map[string]string) string {
	result := s
	for _, val := range secrets {
		if val != "" {
			result = strings.ReplaceAll(result, val, MaskValue(val))
		}
	}
	return result
}
