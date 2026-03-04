package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckConstraint(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		version    string
		satisfied  bool
		wantErr    bool
	}{
		// Exact match
		{"exact match", "1.2.3", "1.2.3", true, false},
		{"exact mismatch", "1.2.3", "1.2.4", false, false},

		// Greater than or equal
		{"gte satisfied", ">=18.0.0", "20.11.0", true, false},
		{"gte exact", ">=18.0.0", "18.0.0", true, false},
		{"gte not satisfied", ">=18.0.0", "16.0.0", false, false},

		// Caret (^) - compatible with version
		{"caret satisfied", "^1.2.3", "1.9.0", true, false},
		{"caret exact", "^1.2.3", "1.2.3", true, false},
		{"caret too low", "^1.2.3", "1.2.2", false, false},
		{"caret major bump", "^1.2.3", "2.0.0", false, false},
		{"caret two-part", "^3.10", "3.12.1", true, false},
		{"caret two-part too low", "^3.10", "3.9.0", false, false},
		{"caret two-part major bump", "^3.10", "4.0.0", false, false},

		// Tilde (~) - patch level changes
		{"tilde satisfied", "~1.2.3", "1.2.9", true, false},
		{"tilde exact", "~1.2.3", "1.2.3", true, false},
		{"tilde minor bump", "~1.2.3", "1.3.0", false, false},
		{"tilde two-part", "~1.2", "1.2.5", true, false},

		// Range
		{"range satisfied", ">=1.0.0, <2.0.0", "1.5.0", true, false},
		{"range too high", ">=1.0.0, <2.0.0", "2.0.0", false, false},
		{"range too low", ">=1.0.0, <2.0.0", "0.9.0", false, false},

		// Star / empty (any)
		{"star any", "*", "99.99.99", true, false},
		{"empty any", "", "1.0.0", true, false},

		// Leading v
		{"leading v version", ">=1.0.0", "v1.2.3", true, false},

		// Two-part versions
		{"two-part version", ">=3.10", "3.12.1", true, false},
		{"two-part version fail", ">=3.10", "3.9.0", false, false},

		// Single-part version
		{"single-part", ">=3", "3.12.1", true, false},

		// Invalid inputs
		{"invalid version", ">=1.0.0", "notaversion", false, true},
		{"invalid constraint", ">>>1.0.0", "1.0.0", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := CheckConstraint(tt.constraint, tt.version)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.satisfied, ok,
				"CheckConstraint(%q, %q) = %v, want %v",
				tt.constraint, tt.version, ok, tt.satisfied)
		})
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"1.2.3", "1.2.3"},
		{"v1.2.3", "1.2.3"},
		{"1.2", "1.2.0"},
		{"1", "1.0.0"},
		{"  v20.11.0 ", "20.11.0"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expect, normalizeVersion(tt.input))
		})
	}
}

func TestNormalizeConstraint(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{">=18.0.0", ">=18.0.0"},
		{"^3.10", "^3.10.0"},
		{"~1.2", "~1.2.0"},
		{">=1.0.0, <2.0.0", ">=1.0.0, <2.0.0"},
		{"*", "*"},
		{"", ""},
		{"^3", "^3.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expect, normalizeConstraint(tt.input))
		})
	}
}
