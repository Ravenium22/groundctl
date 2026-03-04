package drift

import (
	"testing"

	"github.com/Ravenium22/groundctl/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompare_AllOK(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=18.0.0", Severity: model.SeverityRequired},
			{Name: "go", Version: "^1.21", Severity: model.SeverityRequired},
		},
	}
	detected := []model.DetectedTool{
		{Name: "node", Version: "20.11.0", Found: true},
		{Name: "go", Version: "1.22.0", Found: true},
	}

	report := Compare(cfg, detected)
	assert.Equal(t, 2, report.Summary.Total)
	assert.Equal(t, 2, report.Summary.OK)
	assert.Equal(t, 0, report.Summary.Warnings)
	assert.Equal(t, 0, report.Summary.Errors)

	for _, item := range report.Items {
		assert.Equal(t, model.DriftOK, item.Status)
	}
}

func TestCompare_MissingRequired(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "docker", Version: ">=24.0.0", Severity: model.SeverityRequired},
		},
	}
	detected := []model.DetectedTool{
		{Name: "docker", Found: false},
	}

	report := Compare(cfg, detected)
	require.Len(t, report.Items, 1)
	assert.Equal(t, model.DriftError, report.Items[0].Status)
	assert.Equal(t, "not installed", report.Items[0].Message)
	assert.Equal(t, 1, report.Summary.Errors)
}

func TestCompare_MissingRecommended(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "wget", Version: "*", Severity: model.SeverityRecommended},
		},
	}
	detected := []model.DetectedTool{
		{Name: "wget", Found: false},
	}

	report := Compare(cfg, detected)
	require.Len(t, report.Items, 1)
	assert.Equal(t, model.DriftWarning, report.Items[0].Status)
	assert.Equal(t, 0, report.Summary.Errors)
	assert.Equal(t, 1, report.Summary.Warnings)
}

func TestCompare_VersionDrift(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=20.0.0", Severity: model.SeverityRequired},
		},
	}
	detected := []model.DetectedTool{
		{Name: "node", Version: "18.19.0", Found: true},
	}

	report := Compare(cfg, detected)
	require.Len(t, report.Items, 1)
	assert.Equal(t, model.DriftError, report.Items[0].Status)
	assert.Contains(t, report.Items[0].Message, "does not satisfy")
	assert.Equal(t, "18.19.0", report.Items[0].Actual)
	assert.Equal(t, ">=20.0.0", report.Items[0].Expected)
}

func TestCompare_VersionDriftRecommended(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "python", Version: ">=3.12", Severity: model.SeverityRecommended},
		},
	}
	detected := []model.DetectedTool{
		{Name: "python", Version: "3.10.0", Found: true},
	}

	report := Compare(cfg, detected)
	require.Len(t, report.Items, 1)
	assert.Equal(t, model.DriftWarning, report.Items[0].Status)
}

func TestCompare_NoVersionConstraint(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "git", Severity: model.SeverityRequired},
		},
	}
	detected := []model.DetectedTool{
		{Name: "git", Version: "2.43.0", Found: true},
	}

	report := Compare(cfg, detected)
	require.Len(t, report.Items, 1)
	assert.Equal(t, model.DriftOK, report.Items[0].Status)
}

func TestCompare_StarConstraint(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "make", Version: "*", Severity: model.SeverityRequired},
		},
	}
	detected := []model.DetectedTool{
		{Name: "make", Version: "4.4.1", Found: true},
	}

	report := Compare(cfg, detected)
	require.Len(t, report.Items, 1)
	assert.Equal(t, model.DriftOK, report.Items[0].Status)
}

func TestCompare_ToolNotInDetected(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "terraform", Version: ">=1.0.0", Severity: model.SeverityRequired},
		},
	}
	// terraform not in detected list at all
	detected := []model.DetectedTool{}

	report := Compare(cfg, detected)
	require.Len(t, report.Items, 1)
	assert.Equal(t, model.DriftError, report.Items[0].Status)
	assert.Equal(t, "not installed", report.Items[0].Message)
}

func TestCompare_DefaultSeverity(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "curl"}, // no severity set, should default to required
		},
	}
	detected := []model.DetectedTool{}

	report := Compare(cfg, detected)
	require.Len(t, report.Items, 1)
	assert.Equal(t, model.DriftError, report.Items[0].Status)
}

func TestCompare_MixedResults(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=18.0.0", Severity: model.SeverityRequired},
			{Name: "python", Version: ">=3.12", Severity: model.SeverityRecommended},
			{Name: "docker", Version: ">=24.0.0", Severity: model.SeverityRequired},
			{Name: "git"},
		},
	}
	detected := []model.DetectedTool{
		{Name: "node", Version: "20.11.0", Found: true},
		{Name: "python", Version: "3.10.0", Found: true},
		{Name: "docker", Found: false},
		{Name: "git", Version: "2.43.0", Found: true},
	}

	report := Compare(cfg, detected)
	assert.Equal(t, 4, report.Summary.Total)
	assert.Equal(t, 2, report.Summary.OK)       // node + git
	assert.Equal(t, 1, report.Summary.Warnings)  // python (recommended, version drift)
	assert.Equal(t, 1, report.Summary.Errors)     // docker (required, missing)
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name     string
		summary  model.Summary
		expected int
	}{
		{"all clean", model.Summary{OK: 5}, 0},
		{"warnings only", model.Summary{OK: 3, Warnings: 2}, 1},
		{"errors present", model.Summary{OK: 2, Warnings: 1, Errors: 1}, 2},
		{"errors only", model.Summary{Errors: 3}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &model.DriftReport{Summary: tt.summary}
			assert.Equal(t, tt.expected, ExitCode(report))
		})
	}
}
