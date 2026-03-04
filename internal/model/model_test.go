package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNewSnapshot(t *testing.T) {
	tools := []DetectedTool{
		{Name: "node", Version: "20.11.0", Found: true, Path: "/usr/bin/node"},
		{Name: "python", Version: "3.12.1", Found: true},
	}

	snap := NewSnapshot("testhost", "linux", "amd64", tools)
	assert.Equal(t, "testhost", snap.Hostname)
	assert.Equal(t, "linux", snap.OS)
	assert.Equal(t, "amd64", snap.Arch)
	assert.Len(t, snap.Tools, 2)
	assert.NotEmpty(t, snap.Timestamp)
}

func TestSnapshotJSON(t *testing.T) {
	snap := NewSnapshot("host", "darwin", "arm64", []DetectedTool{
		{Name: "go", Version: "1.22.0", Found: true},
	})

	data, err := json.Marshal(snap)
	require.NoError(t, err)

	var decoded Snapshot
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, snap.Hostname, decoded.Hostname)
	assert.Len(t, decoded.Tools, 1)
	assert.Equal(t, "go", decoded.Tools[0].Name)
}

func TestGroundConfigYAML(t *testing.T) {
	cfg := GroundConfig{
		Name:        "my-project",
		Description: "Test config",
		Tools: []ToolSpec{
			{Name: "node", Version: ">=18.0.0", Severity: SeverityRequired},
			{Name: "python", Version: "^3.10", Severity: SeverityRecommended},
		},
	}

	data, err := yaml.Marshal(&cfg)
	require.NoError(t, err)

	var decoded GroundConfig
	err = yaml.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "my-project", decoded.Name)
	assert.Len(t, decoded.Tools, 2)
	assert.Equal(t, SeverityRequired, decoded.Tools[0].Severity)
	assert.Equal(t, ">=18.0.0", decoded.Tools[0].Version)
}

func TestDriftReport(t *testing.T) {
	report := DriftReport{
		Timestamp: "2024-01-01T00:00:00Z",
		Items: []DriftItem{
			{Tool: "node", Status: DriftOK, Expected: ">=18.0.0", Actual: "20.11.0"},
			{Tool: "python", Status: DriftWarning, Expected: "^3.12", Actual: "3.10.0", Message: "version drift"},
			{Tool: "docker", Status: DriftError, Message: "not installed"},
		},
		Summary: Summary{Total: 3, OK: 1, Warnings: 1, Errors: 1},
	}

	data, err := json.Marshal(report)
	require.NoError(t, err)

	var decoded DriftReport
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, 3, decoded.Summary.Total)
	assert.Equal(t, DriftOK, decoded.Items[0].Status)
	assert.Equal(t, DriftError, decoded.Items[2].Status)
}
