package team

import (
	"path/filepath"
	"testing"

	"github.com/groundctl/groundctl/internal/config"
	"github.com/groundctl/groundctl/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullLocal(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	// Create a source config
	src := filepath.Join(srcDir, ".ground.yaml")
	cfg := &model.GroundConfig{
		Name: "team-config",
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=18.0.0", Severity: model.SeverityRequired},
			{Name: "go", Version: ">=1.21", Severity: model.SeverityRequired},
		},
	}
	require.NoError(t, config.Save(src, cfg))

	// Pull
	dest := filepath.Join(destDir, ".ground.yaml")
	pulled, err := Pull(src, dest)
	require.NoError(t, err)
	assert.Equal(t, "team-config", pulled.Name)
	assert.Len(t, pulled.Tools, 2)

	// Verify dest file exists
	assert.True(t, config.Exists(dest))
}

func TestPullLocalNotFound(t *testing.T) {
	_, err := Pull("/nonexistent/path/.ground.yaml", "/tmp/test.yaml")
	assert.Error(t, err)
}

func TestIsLocalPath(t *testing.T) {
	assert.True(t, isLocalPath("/absolute/path"))
	assert.True(t, isLocalPath("./relative/path"))
	assert.True(t, isLocalPath("../parent/path"))
	assert.True(t, isLocalPath("C:\\Windows\\path"))
	assert.False(t, isLocalPath("github.com/org/repo"))
}

func TestIsDirectURL(t *testing.T) {
	assert.True(t, isDirectURL("https://example.com/config.yaml"))
	assert.True(t, isDirectURL("http://example.com/config.yml"))
	assert.False(t, isDirectURL("https://github.com/org/repo"))
	assert.False(t, isDirectURL("github.com/org/repo"))
}

func TestNormalizeGitURL(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"github.com/org/repo", "https://github.com/org/repo"},
		{"https://github.com/org/repo", "https://github.com/org/repo"},
		{"https://github.com/org/repo/", "https://github.com/org/repo"},
		{"git@github.com:org/repo.git", "git@github.com:org/repo.git"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expect, normalizeGitURL(tt.input))
		})
	}
}

func TestPullGitNoRepo(t *testing.T) {
	dest := filepath.Join(t.TempDir(), ".ground.yaml")
	_, err := Pull("https://github.com/nonexistent-org-12345/nonexistent-repo-67890", dest)
	assert.Error(t, err)
}

func TestPushNoConfig(t *testing.T) {
	// Push should fail when git add fails (nonexistent file)
	err := Push("nonexistent-file-that-does-not-exist.yaml", "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "git add failed")
}
