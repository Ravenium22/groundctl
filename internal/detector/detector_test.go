package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"node version", "v20.11.0", "20.11.0"},
		{"npm version", "10.2.4", "10.2.4"},
		{"python version", "Python 3.12.1", "3.12.1"},
		{"go version", "go version go1.22.0 windows/amd64", "1.22.0"},
		{"git version", "git version 2.43.0.windows.1", "2.43.0"},
		{"docker version", "Docker version 24.0.7, build afdd53b", "24.0.7"},
		{"rust version", "rustc 1.75.0 (82e1608df 2023-12-21)", "1.75.0"},
		{"ruby version", "ruby 3.3.0 (2023-12-25 revision 5124f9ac75)", "3.3.0"},
		{"java version", "openjdk version \"21.0.1\" 2023-10-17", "21.0.1"},
		{"gh version", "gh version 2.42.0 (2024-01-15)", "2.42.0"},
		{"terraform version", "Terraform v1.7.0", "1.7.0"},
		{"kubectl version", "Client Version: v1.29.0", "1.29.0"},
		{"make version", "GNU Make 4.4.1", "4.4.1"},
		{"cargo version", "cargo 1.75.0 (1d8b05cdd 2023-11-20)", "1.75.0"},
		{"curl version", "curl 8.4.0 (x86_64-pc-linux-gnu)", "8.4.0"},
		{"no version", "some random output", ""},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseVersion(tt.input)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestListKnownTools(t *testing.T) {
	tools := ListKnownTools()
	require.True(t, len(tools) >= 15, "should have at least 15 detectable tools, got %d", len(tools))

	expected := []string{"node", "npm", "python", "go", "git", "docker", "gh"}
	for _, name := range expected {
		assert.Contains(t, tools, name, "should include %s", name)
	}
}

func TestDetectAll(t *testing.T) {
	results := DetectAll()
	require.NotEmpty(t, results, "should detect at least some tools")

	// At minimum, git should be available in most CI/dev environments
	var foundAny bool
	for _, r := range results {
		if r.Found {
			foundAny = true
			assert.NotEmpty(t, r.Version, "found tool %s should have a version", r.Name)
		}
	}
	// We don't assert foundAny is true because minimal environments may have nothing
	_ = foundAny
}

func TestDetectByNames(t *testing.T) {
	results := DetectByNames([]string{"git"})
	require.Len(t, results, 1, "should return exactly 1 result for 1 name")
	assert.Equal(t, "git", results[0].Name)
}

func TestDetectByNamesEmpty(t *testing.T) {
	results := DetectByNames([]string{})
	assert.Empty(t, results)
}

func TestDetectByNamesUnknown(t *testing.T) {
	results := DetectByNames([]string{"nonexistenttool12345"})
	assert.Empty(t, results, "unknown tool names should return empty")
}
