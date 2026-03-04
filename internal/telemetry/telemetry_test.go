package telemetry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsEnabledDefault(t *testing.T) {
	// Reset state
	mu.Lock()
	configured = false
	enabled = false
	mu.Unlock()

	// With no config file and no env var, should be disabled
	t.Setenv("GROUND_TELEMETRY", "")
	assert.False(t, IsEnabled())
}

func TestIsEnabledEnvOff(t *testing.T) {
	mu.Lock()
	configured = false
	enabled = false
	mu.Unlock()

	t.Setenv("GROUND_TELEMETRY", "off")
	assert.False(t, IsEnabled())
}

func TestIsEnabledEnvFalse(t *testing.T) {
	mu.Lock()
	configured = false
	enabled = false
	mu.Unlock()

	t.Setenv("GROUND_TELEMETRY", "false")
	assert.False(t, IsEnabled())
}

func TestSetEnabledAndRecord(t *testing.T) {
	// Use a temp home dir
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)
	t.Setenv("GROUND_TELEMETRY", "")

	mu.Lock()
	configured = false
	enabled = false
	mu.Unlock()

	// Enable telemetry
	err := SetEnabled(true)
	require.NoError(t, err)

	// Verify config file exists
	cfgPath := filepath.Join(tmpDir, ".groundctl", "telemetry.json")
	_, err = os.Stat(cfgPath)
	assert.NoError(t, err)

	assert.True(t, IsEnabled())

	// Record an event
	Record(Event{
		Command:    "check",
		ToolCount:  5,
		DurationMs: 150,
		ExitCode:   0,
	})

	// Check event was stored
	assert.Equal(t, 1, EventCount())

	// Record another
	Record(Event{
		Command:    "fix",
		ToolCount:  3,
		DurationMs: 2000,
		ExitCode:   0,
	})
	assert.Equal(t, 2, EventCount())

	// Clear events
	err = ClearEvents()
	assert.NoError(t, err)
	assert.Equal(t, 0, EventCount())
}

func TestSetEnabledFalse(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)
	t.Setenv("GROUND_TELEMETRY", "")

	mu.Lock()
	configured = false
	enabled = false
	mu.Unlock()

	err := SetEnabled(false)
	require.NoError(t, err)
	assert.False(t, IsEnabled())

	// Record should be a no-op
	Record(Event{Command: "test"})
	assert.Equal(t, 0, EventCount())
}

func TestRecordDisabled(t *testing.T) {
	mu.Lock()
	configured = true
	enabled = false
	mu.Unlock()

	// Should not panic or error
	Record(Event{Command: "test"})
}
