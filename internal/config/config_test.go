package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/groundctl/groundctl/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".ground.yaml")

	cfg := &model.GroundConfig{
		Name:        "test-project",
		Description: "A test config",
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=18.0.0", Severity: model.SeverityRequired},
			{Name: "go", Version: ">=1.21.0", Severity: model.SeverityRequired},
		},
	}

	err := Save(path, cfg)
	require.NoError(t, err)

	loaded, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "test-project", loaded.Name)
	assert.Len(t, loaded.Tools, 2)
	assert.Equal(t, "node", loaded.Tools[0].Name)
	assert.Equal(t, ">=18.0.0", loaded.Tools[0].Version)
}

func TestLoadNonExistent(t *testing.T) {
	_, err := Load("/nonexistent/.ground.yaml")
	assert.Error(t, err)
}

func TestExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".ground.yaml")

	assert.False(t, Exists(path))

	err := os.WriteFile(path, []byte("test"), 0644)
	require.NoError(t, err)

	assert.True(t, Exists(path))
}

func TestSaveContainsHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".ground.yaml")

	cfg := &model.GroundConfig{Name: "test", Tools: []model.ToolSpec{}}
	require.NoError(t, Save(path, cfg))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "# groundctl configuration")
}
