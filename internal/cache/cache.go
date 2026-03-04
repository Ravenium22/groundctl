package cache

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Ravenium22/groundctl/internal/model"
)

// DefaultTTL is the default cache lifetime.
var DefaultTTL = 5 * time.Minute

// Entry represents a cached detection result.
type Entry struct {
	Tool      model.DetectedTool `json:"tool"`
	BinaryMod string             `json:"binary_mod,omitempty"` // mtime of binary
	CachedAt  time.Time          `json:"cached_at"`
}

// Store holds cached detection entries.
type Store struct {
	Entries map[string]Entry `json:"entries"`
	path    string
	ttl     time.Duration
}

// New creates a new cache store at the default location.
func New() *Store {
	home, err := os.UserHomeDir()
	if err != nil {
		return &Store{Entries: make(map[string]Entry), ttl: DefaultTTL}
	}
	return &Store{
		Entries: make(map[string]Entry),
		path:    filepath.Join(home, ".groundctl", "cache", "detect.json"),
		ttl:     DefaultTTL,
	}
}

// NewWithPath creates a cache store at a specific path.
func NewWithPath(path string, ttl time.Duration) *Store {
	return &Store{
		Entries: make(map[string]Entry),
		path:    path,
		ttl:     ttl,
	}
}

// Load reads cached entries from disk.
func (s *Store) Load() error {
	if s.path == "" {
		return nil
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.Entries)
}

// Save writes cached entries to disk.
func (s *Store) Save() error {
	if s.path == "" {
		return nil
	}
	_ = os.MkdirAll(filepath.Dir(s.path), 0755)
	data, err := json.MarshalIndent(s.Entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// Get retrieves a cached entry if it exists and is fresh.
func (s *Store) Get(toolName string) (model.DetectedTool, bool) {
	entry, ok := s.Entries[toolName]
	if !ok {
		return model.DetectedTool{}, false
	}

	// Check TTL
	if time.Since(entry.CachedAt) > s.ttl {
		delete(s.Entries, toolName)
		return model.DetectedTool{}, false
	}

	// Check if binary has changed (mtime invalidation)
	if entry.Tool.Path != "" {
		currentMod := binaryMtime(entry.Tool.Path)
		if currentMod != entry.BinaryMod {
			delete(s.Entries, toolName)
			return model.DetectedTool{}, false
		}
	}

	return entry.Tool, true
}

// Put stores a detection result in the cache.
func (s *Store) Put(tool model.DetectedTool) {
	entry := Entry{
		Tool:     tool,
		CachedAt: time.Now(),
	}
	if tool.Path != "" {
		entry.BinaryMod = binaryMtime(tool.Path)
	}
	s.Entries[tool.Name] = entry
}

// Clear removes all cached entries.
func (s *Store) Clear() {
	s.Entries = make(map[string]Entry)
	if s.path != "" {
		_ = os.Remove(s.path)
	}
}

// Size returns the number of cached entries.
func (s *Store) Size() int {
	return len(s.Entries)
}

func binaryMtime(path string) string {
	// Resolve symlinks first
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		resolved = path
	}
	info, err := os.Stat(resolved)
	if err != nil {
		// If the binary path came from LookPath, try that too
		if lp, err2 := exec.LookPath(filepath.Base(resolved)); err2 == nil {
			if info2, err3 := os.Stat(lp); err3 == nil {
				return info2.ModTime().UTC().Format(time.RFC3339Nano)
			}
		}
		return ""
	}
	return info.ModTime().UTC().Format(time.RFC3339Nano)
}
