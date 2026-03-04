package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/Ravenium22/groundctl/internal/detector"
	"github.com/Ravenium22/groundctl/internal/model"
)

// Capture runs all detectors and returns a full machine snapshot.
func Capture() (*model.Snapshot, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	tools := detector.DetectAll()
	snap := model.NewSnapshot(hostname, runtime.GOOS, runtime.GOARCH, tools)
	return snap, nil
}

// CaptureForTools runs detectors for specific tool names only.
func CaptureForTools(names []string) (*model.Snapshot, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	tools := detector.DetectByNames(names)
	snap := model.NewSnapshot(hostname, runtime.GOOS, runtime.GOARCH, tools)
	return snap, nil
}

// ToJSON marshals a snapshot to indented JSON.
func ToJSON(snap *model.Snapshot) ([]byte, error) {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("could not marshal snapshot: %w", err)
	}
	return data, nil
}
