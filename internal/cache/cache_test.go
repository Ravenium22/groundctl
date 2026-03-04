package cache

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/groundctl/groundctl/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	s := New()
	assert.NotNil(t, s)
	assert.Equal(t, 0, s.Size())
}

func TestPutAndGet(t *testing.T) {
	s := NewWithPath("", 5*time.Minute)

	tool := model.DetectedTool{
		Name:    "git",
		Version: "2.43.0",
		Path:    "/usr/bin/git",
		Found:   true,
	}

	s.Put(tool)
	assert.Equal(t, 1, s.Size())

	got, ok := s.Get("git")
	assert.True(t, ok)
	assert.Equal(t, "git", got.Name)
	assert.Equal(t, "2.43.0", got.Version)
}

func TestGetMiss(t *testing.T) {
	s := NewWithPath("", 5*time.Minute)
	_, ok := s.Get("nonexistent")
	assert.False(t, ok)
}

func TestTTLExpiry(t *testing.T) {
	s := NewWithPath("", 1*time.Millisecond)

	tool := model.DetectedTool{
		Name:    "node",
		Version: "20.0.0",
		Found:   true,
	}
	s.Put(tool)

	// Wait for TTL to expire
	time.Sleep(5 * time.Millisecond)

	_, ok := s.Get("node")
	assert.False(t, ok)
}

func TestClear(t *testing.T) {
	s := NewWithPath("", 5*time.Minute)
	s.Put(model.DetectedTool{Name: "a", Found: true})
	s.Put(model.DetectedTool{Name: "b", Found: true})
	assert.Equal(t, 2, s.Size())

	s.Clear()
	assert.Equal(t, 0, s.Size())
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "cache.json")

	// Save
	s1 := NewWithPath(path, 5*time.Minute)
	s1.Put(model.DetectedTool{Name: "go", Version: "1.22.0", Found: true})
	s1.Put(model.DetectedTool{Name: "node", Version: "20.10.0", Found: true})
	err := s1.Save()
	require.NoError(t, err)

	// Load into new store
	s2 := NewWithPath(path, 5*time.Minute)
	err = s2.Load()
	require.NoError(t, err)

	assert.Equal(t, 2, s2.Size())
	got, ok := s2.Get("go")
	assert.True(t, ok)
	assert.Equal(t, "1.22.0", got.Version)
}

func TestLoadMissingFile(t *testing.T) {
	s := NewWithPath("/nonexistent/path/cache.json", 5*time.Minute)
	err := s.Load()
	assert.Error(t, err)
}

func TestPutNotFoundTool(t *testing.T) {
	s := NewWithPath("", 5*time.Minute)
	tool := model.DetectedTool{Name: "missing", Found: false, Error: "not found"}
	s.Put(tool)

	got, ok := s.Get("missing")
	assert.True(t, ok)
	assert.False(t, got.Found)
}
