package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Ravenium22/groundctl/internal/config"
	"github.com/Ravenium22/groundctl/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDir overrides the profile directory for testing.
func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
	profileDir := filepath.Join(dir, ".groundctl", "profiles")
	require.NoError(t, os.MkdirAll(profileDir, 0755))
	return dir
}

func TestSaveAndLoad(t *testing.T) {
	setupTestDir(t)

	cfg := &model.GroundConfig{
		Name: "test-profile",
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=18.0.0", Severity: model.SeverityRequired},
			{Name: "go", Version: ">=1.21", Severity: model.SeverityRequired},
		},
	}

	err := Save("work", cfg)
	require.NoError(t, err)

	loaded, err := Load("work")
	require.NoError(t, err)
	assert.Equal(t, "test-profile", loaded.Name)
	assert.Len(t, loaded.Tools, 2)
}

func TestList(t *testing.T) {
	setupTestDir(t)

	// Save a couple of profiles
	cfg := &model.GroundConfig{Name: "p1", Tools: []model.ToolSpec{{Name: "git"}}}
	require.NoError(t, Save("work", cfg))
	require.NoError(t, Save("personal", cfg))

	names, err := List()
	require.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "work")
	assert.Contains(t, names, "personal")
}

func TestListEmpty(t *testing.T) {
	setupTestDir(t)
	names, err := List()
	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestExists(t *testing.T) {
	setupTestDir(t)

	assert.False(t, Exists("nonexistent"))

	cfg := &model.GroundConfig{Name: "test", Tools: []model.ToolSpec{{Name: "git"}}}
	require.NoError(t, Save("test", cfg))

	assert.True(t, Exists("test"))
}

func TestDelete(t *testing.T) {
	setupTestDir(t)

	cfg := &model.GroundConfig{Name: "test", Tools: []model.ToolSpec{{Name: "git"}}}
	require.NoError(t, Save("temp", cfg))
	assert.True(t, Exists("temp"))

	require.NoError(t, Delete("temp"))
	assert.False(t, Exists("temp"))
}

func TestDeleteNonExistent(t *testing.T) {
	setupTestDir(t)
	err := Delete("nonexistent")
	assert.NoError(t, err)
}

func TestActiveProfile(t *testing.T) {
	setupTestDir(t)

	assert.Equal(t, "", GetActive())

	require.NoError(t, SetActive("work"))
	assert.Equal(t, "work", GetActive())

	require.NoError(t, SetActive("personal"))
	assert.Equal(t, "personal", GetActive())

	require.NoError(t, ClearActive())
	assert.Equal(t, "", GetActive())
}

func TestResolve_NoInheritance(t *testing.T) {
	setupTestDir(t)

	cfg := &model.GroundConfig{
		Name:  "base",
		Tools: []model.ToolSpec{{Name: "git"}, {Name: "node", Version: ">=18"}},
	}
	require.NoError(t, Save("base", cfg))

	resolved, err := Resolve("base")
	require.NoError(t, err)
	assert.Len(t, resolved.Tools, 2)
}

func TestResolve_WithInheritance(t *testing.T) {
	setupTestDir(t)

	// Parent profile
	parent := &model.GroundConfig{
		Name: "team-base",
		Tools: []model.ToolSpec{
			{Name: "git", Version: ">=2.40"},
			{Name: "node", Version: ">=18"},
			{Name: "docker", Version: ">=24"},
		},
	}
	require.NoError(t, Save("team-base", parent))

	// Child profile extends parent
	child := &model.GroundConfig{
		Name:    "my-team",
		Extends: "team-base",
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=20"}, // override parent
			{Name: "go", Version: ">=1.21"}, // add new
		},
	}
	require.NoError(t, Save("my-team", child))

	resolved, err := Resolve("my-team")
	require.NoError(t, err)

	// Should have: node (child override), go (child), git (parent), docker (parent) = 4
	assert.Len(t, resolved.Tools, 4)

	// node should use child version (>=20), not parent (>=18)
	for _, tool := range resolved.Tools {
		if tool.Name == "node" {
			assert.Equal(t, ">=20", tool.Version)
		}
	}
}

func TestResolve_ChainedInheritance(t *testing.T) {
	setupTestDir(t)

	grandparent := &model.GroundConfig{
		Name:  "org-base",
		Tools: []model.ToolSpec{{Name: "git"}, {Name: "curl"}},
	}
	require.NoError(t, Save("org-base", grandparent))

	parent := &model.GroundConfig{
		Name:    "team-base",
		Extends: "org-base",
		Tools:   []model.ToolSpec{{Name: "node", Version: ">=18"}},
	}
	require.NoError(t, Save("team-base", parent))

	child := &model.GroundConfig{
		Name:    "my-env",
		Extends: "team-base",
		Tools:   []model.ToolSpec{{Name: "go", Version: ">=1.21"}},
	}
	require.NoError(t, Save("my-env", child))

	resolved, err := Resolve("my-env")
	require.NoError(t, err)
	// go (child) + node (parent) + git, curl (grandparent) = 4
	assert.Len(t, resolved.Tools, 4)
}

func TestResolve_MissingParent(t *testing.T) {
	setupTestDir(t)

	child := &model.GroundConfig{
		Name:    "orphan",
		Extends: "nonexistent-parent",
		Tools:   []model.ToolSpec{{Name: "git"}},
	}
	require.NoError(t, Save("orphan", child))

	_, err := Resolve("orphan")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent-parent")
}

func TestMergeTools(t *testing.T) {
	child := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=20"},
		},
	}
	parent := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=18"},
			{Name: "git", Version: ">=2.40"},
		},
	}

	child.MergeTools(parent)
	assert.Len(t, child.Tools, 2)

	// node should remain child version
	for _, tool := range child.Tools {
		if tool.Name == "node" {
			assert.Equal(t, ">=20", tool.Version)
		}
	}
}

func TestMergeToolsNilParent(t *testing.T) {
	child := &model.GroundConfig{
		Tools: []model.ToolSpec{{Name: "node"}},
	}
	child.MergeTools(nil) // should not panic
	assert.Len(t, child.Tools, 1)
}

// Helper to ensure test profiles use config.Save (which we know works)
func init() {
	_ = config.DefaultConfigFile
}
