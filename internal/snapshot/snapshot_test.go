package snapshot

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/groundctl/groundctl/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapture(t *testing.T) {
	snap, err := Capture()
	require.NoError(t, err)
	require.NotNil(t, snap)

	assert.NotEmpty(t, snap.Timestamp)
	assert.NotEmpty(t, snap.Hostname)
	assert.Equal(t, runtime.GOOS, snap.OS)
	assert.Equal(t, runtime.GOARCH, snap.Arch)
	assert.NotEmpty(t, snap.Tools, "should detect at least some tools")
}

func TestCaptureForTools(t *testing.T) {
	snap, err := CaptureForTools([]string{"git"})
	require.NoError(t, err)
	require.NotNil(t, snap)
	assert.Len(t, snap.Tools, 1)
	assert.Equal(t, "git", snap.Tools[0].Name)
}

func TestToJSON(t *testing.T) {
	snap := &model.Snapshot{
		Timestamp: "2024-01-01T00:00:00Z",
		Hostname:  "test",
		OS:        "linux",
		Arch:      "amd64",
		Tools: []model.DetectedTool{
			{Name: "go", Version: "1.22.0", Found: true},
		},
	}

	data, err := ToJSON(snap)
	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &parsed))
	assert.Equal(t, "test", parsed["hostname"])
}
