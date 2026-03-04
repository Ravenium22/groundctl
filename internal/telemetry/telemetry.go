package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var eventCounter atomic.Int64

// Event represents a single telemetry event.
type Event struct {
	Command   string `json:"command"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	ToolCount int    `json:"tool_count,omitempty"`
	DurationMs int64 `json:"duration_ms"`
	ExitCode  int    `json:"exit_code"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version,omitempty"`
}

// Config holds telemetry opt-in state.
type Config struct {
	Enabled   bool   `json:"enabled"`
	Endpoint  string `json:"endpoint,omitempty"`
	UpdatedAt string `json:"updated_at"`
}

var (
	enabled    bool
	configured bool
	mu         sync.Mutex
)

// IsEnabled returns whether telemetry is active.
func IsEnabled() bool {
	mu.Lock()
	defer mu.Unlock()

	if configured {
		return enabled
	}
	configured = true

	// Environment variable override (highest priority)
	if env := os.Getenv("GROUND_TELEMETRY"); env == "off" || env == "false" || env == "0" {
		enabled = false
		return false
	}

	// Check config file
	cfg, err := loadConfig()
	if err != nil || !cfg.Enabled {
		enabled = false
		return false
	}

	enabled = true
	return true
}

// SetEnabled updates the telemetry opt-in state.
func SetEnabled(val bool) error {
	cfg := Config{
		Enabled:   val,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	mu.Lock()
	enabled = val
	configured = true
	mu.Unlock()

	return saveConfig(cfg)
}

// Record stores a telemetry event locally.
// Events are stored in ~/.groundctl/telemetry/events/ as JSON files.
// In a production build, these would be batched and sent to the endpoint.
func Record(evt Event) {
	if !IsEnabled() {
		return
	}

	evt.OS = runtime.GOOS
	evt.Arch = runtime.GOARCH
	evt.Timestamp = time.Now().UTC().Format(time.RFC3339)

	eventsDir, err := eventsPath()
	if err != nil {
		return
	}
	_ = os.MkdirAll(eventsDir, 0755)

	seq := eventCounter.Add(1)
	filename := fmt.Sprintf("%d-%d.json", time.Now().UnixNano(), seq)
	data, err := json.Marshal(evt)
	if err != nil {
		return
	}

	_ = os.WriteFile(filepath.Join(eventsDir, filename), data, 0644)
}

// EventCount returns the number of stored telemetry events.
func EventCount() int {
	dir, err := eventsPath()
	if err != nil {
		return 0
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	return count
}

// ClearEvents removes all stored telemetry events.
func ClearEvents() error {
	dir, err := eventsPath()
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".groundctl", "telemetry.json"), nil
}

func eventsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".groundctl", "telemetry", "events"), nil
}

func loadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
