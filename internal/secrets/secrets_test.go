package secrets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Parser tests ---

func TestParseRef(t *testing.T) {
	tests := []struct {
		input   string
		want    *SecretRef
	}{
		{"${op://vault/item/field}", &SecretRef{Raw: "${op://vault/item/field}", Backend: "op", Path: "vault/item/field"}},
		{"${vault://secret/db#password}", &SecretRef{Raw: "${vault://secret/db#password}", Backend: "vault", Path: "secret/db#password"}},
		{"${env://DATABASE_URL}", &SecretRef{Raw: "${env://DATABASE_URL}", Backend: "env", Path: "DATABASE_URL"}},
		{"${keychain://myapp/token}", &SecretRef{Raw: "${keychain://myapp/token}", Backend: "keychain", Path: "myapp/token"}},
		{"not a ref", nil},
		{"${invalid}", nil},
		{"${://no-backend}", nil},
		{"", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseRef(tt.input)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tt.want.Raw, got.Raw)
				assert.Equal(t, tt.want.Backend, got.Backend)
				assert.Equal(t, tt.want.Path, got.Path)
			}
		})
	}
}

func TestParseRefs(t *testing.T) {
	input := "host=${vault://db/host}:${vault://db/port}"
	refs := ParseRefs(input)
	assert.Len(t, refs, 2)
	assert.Equal(t, "vault", refs[0].Backend)
	assert.Equal(t, "db/host", refs[0].Path)
	assert.Equal(t, "vault", refs[1].Backend)
	assert.Equal(t, "db/port", refs[1].Path)
}

func TestParseRefsNoRefs(t *testing.T) {
	refs := ParseRefs("just a plain string")
	assert.Empty(t, refs)
}

// --- Mock backend for testing ---

type mockBackend struct {
	name    string
	secrets map[string]string
	err     error
	checkOK bool
}

func (m *mockBackend) Name() string { return m.name }
func (m *mockBackend) Resolve(path string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	val, ok := m.secrets[path]
	if !ok {
		return "", fmt.Errorf("not found: %s", path)
	}
	return val, nil
}
func (m *mockBackend) Check() error {
	if !m.checkOK {
		return fmt.Errorf("backend %s not available", m.name)
	}
	return nil
}

// --- Registry tests ---

func TestRegistryResolve(t *testing.T) {
	r := NewRegistry()
	r.Register(&mockBackend{
		name:    "mock",
		secrets: map[string]string{"key": "secret-value"},
		checkOK: true,
	})

	ref := SecretRef{Raw: "${mock://key}", Backend: "mock", Path: "key"}
	val, err := r.Resolve(ref)
	require.NoError(t, err)
	assert.Equal(t, "secret-value", val)
}

func TestRegistryResolveUnknownBackend(t *testing.T) {
	r := NewRegistry()
	ref := SecretRef{Raw: "${unknown://key}", Backend: "unknown", Path: "key"}
	_, err := r.Resolve(ref)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown secret backend")
}

func TestRegistryResolveString(t *testing.T) {
	r := NewRegistry()
	r.Register(&mockBackend{
		name: "mock",
		secrets: map[string]string{
			"host": "localhost",
			"port": "5432",
		},
		checkOK: true,
	})

	result, err := r.ResolveString("postgres://${mock://host}:${mock://port}/db")
	require.NoError(t, err)
	assert.Equal(t, "postgres://localhost:5432/db", result)
}

func TestRegistryResolveStringNoRefs(t *testing.T) {
	r := NewRegistry()
	result, err := r.ResolveString("plain string")
	require.NoError(t, err)
	assert.Equal(t, "plain string", result)
}

func TestRegistryResolveStringError(t *testing.T) {
	r := NewRegistry()
	r.Register(&mockBackend{
		name:    "mock",
		secrets: map[string]string{},
		checkOK: true,
	})

	_, err := r.ResolveString("${mock://missing}")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resolving")
}

func TestRegistryCheckRef(t *testing.T) {
	r := NewRegistry()
	r.Register(&mockBackend{
		name:    "mock",
		secrets: map[string]string{"key": "val"},
		checkOK: true,
	})

	ref := SecretRef{Backend: "mock", Path: "key"}
	err := r.CheckRef(ref)
	assert.NoError(t, err)
}

func TestRegistryCheckRefBackendUnavailable(t *testing.T) {
	r := NewRegistry()
	r.Register(&mockBackend{
		name:    "mock",
		secrets: map[string]string{"key": "val"},
		checkOK: false,
	})

	ref := SecretRef{Backend: "mock", Path: "key"}
	err := r.CheckRef(ref)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not available")
}

// --- Env backend test ---

func TestEnvBackend(t *testing.T) {
	b := &EnvBackend{}
	assert.Equal(t, "env", b.Name())
	assert.NoError(t, b.Check())

	t.Setenv("GROUND_TEST_SECRET", "test-value-123")
	val, err := b.Resolve("GROUND_TEST_SECRET")
	require.NoError(t, err)
	assert.Equal(t, "test-value-123", val)
}

func TestEnvBackendMissing(t *testing.T) {
	b := &EnvBackend{}
	_, err := b.Resolve("GROUND_DEFINITELY_NOT_SET_XYZ")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not set")
}

// --- Masking tests ---

func TestMaskValue(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ab", "****"},
		{"abcd", "****"},
		{"abcde", "ab*de"},
		{"secret-value", "se********ue"},
		{"my-long-api-key-12345", "my*****************45"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, MaskValue(tt.input))
		})
	}
}

func TestMaskString(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS":  "supersecret",
		"API_KEY":  "key-12345",
	}

	input := "connecting to db with pass=supersecret and key=key-12345"
	result := MaskString(input, secrets)
	assert.NotContains(t, result, "supersecret")
	assert.NotContains(t, result, "key-12345")
	assert.Contains(t, result, "su*******et")
	assert.Contains(t, result, "ke*****45")
}

func TestMaskStringEmpty(t *testing.T) {
	secrets := map[string]string{"KEY": ""}
	result := MaskString("no secrets here", secrets)
	assert.Equal(t, "no secrets here", result)
}

// --- Default registry test ---

func TestDefaultRegistryHasBackends(t *testing.T) {
	r := DefaultRegistry()

	_, err := r.Get("env")
	assert.NoError(t, err)

	_, err = r.Get("op")
	assert.NoError(t, err)

	_, err = r.Get("vault")
	assert.NoError(t, err)

	_, err = r.Get("keychain")
	assert.NoError(t, err)
}

// --- Op backend name test ---

func TestOpBackendName(t *testing.T) {
	b := &OpBackend{}
	assert.Equal(t, "op", b.Name())
}

// --- Vault backend name test ---

func TestVaultBackendName(t *testing.T) {
	b := &VaultBackend{}
	assert.Equal(t, "vault", b.Name())
}

// --- Keychain backend name test ---

func TestKeychainBackendName(t *testing.T) {
	b := &KeychainBackend{}
	assert.Equal(t, "keychain", b.Name())
}
