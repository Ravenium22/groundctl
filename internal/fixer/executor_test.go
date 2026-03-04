package fixer

import (
	"testing"

	"github.com/groundctl/groundctl/internal/model"
	"github.com/groundctl/groundctl/internal/pkgmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFixPlans_SkipsOK(t *testing.T) {
	report := &model.DriftReport{
		Items: []model.DriftItem{
			{Tool: "node", Status: model.DriftOK, Actual: "20.11.0"},
			{Tool: "go", Status: model.DriftOK, Actual: "1.22.0"},
		},
	}
	managers := []pkgmanager.Manager{{Name: "brew"}}
	plans := BuildFixPlans(report, managers)
	assert.Empty(t, plans, "should have no fix plans when everything is OK")
}

func TestBuildFixPlans_MissingTool(t *testing.T) {
	report := &model.DriftReport{
		Items: []model.DriftItem{
			{Tool: "docker", Status: model.DriftError, Expected: ">=24.0.0", Message: "not installed"},
		},
	}
	managers := []pkgmanager.Manager{{Name: "brew"}}
	plans := BuildFixPlans(report, managers)
	require.Len(t, plans, 1)
	assert.Equal(t, "docker", plans[0].Tool)
	assert.Equal(t, ActionInstall, plans[0].Action)
	assert.Equal(t, "brew", plans[0].Manager)
}

func TestBuildFixPlans_VersionDrift(t *testing.T) {
	report := &model.DriftReport{
		Items: []model.DriftItem{
			{Tool: "node", Status: model.DriftError, Expected: ">=20.0.0", Actual: "18.19.0"},
		},
	}
	managers := []pkgmanager.Manager{{Name: "brew"}}
	plans := BuildFixPlans(report, managers)
	require.Len(t, plans, 1)
	assert.Equal(t, ActionUpgrade, plans[0].Action)
}

func TestBuildFixPlans_WarningIncluded(t *testing.T) {
	report := &model.DriftReport{
		Items: []model.DriftItem{
			{Tool: "python", Status: model.DriftWarning, Expected: ">=3.12", Actual: "3.10.0"},
		},
	}
	managers := []pkgmanager.Manager{{Name: "apt"}}
	plans := BuildFixPlans(report, managers)
	require.Len(t, plans, 1)
	assert.Equal(t, "python", plans[0].Tool)
}

func TestBuildFixPlans_NoPM(t *testing.T) {
	report := &model.DriftReport{
		Items: []model.DriftItem{
			{Tool: "docker", Status: model.DriftError, Expected: ">=24.0.0", Message: "not installed"},
		},
	}
	var managers []pkgmanager.Manager
	plans := BuildFixPlans(report, managers)
	require.Len(t, plans, 1)
	assert.Empty(t, plans[0].Manager)
	assert.NotEmpty(t, plans[0].ManualHint)
}

func TestBuildFixPlans_Mixed(t *testing.T) {
	report := &model.DriftReport{
		Items: []model.DriftItem{
			{Tool: "node", Status: model.DriftOK, Actual: "20.11.0"},
			{Tool: "docker", Status: model.DriftError, Expected: ">=24.0.0"},
			{Tool: "python", Status: model.DriftWarning, Expected: ">=3.12", Actual: "3.10.0"},
			{Tool: "git", Status: model.DriftOK, Actual: "2.43.0"},
		},
	}
	managers := []pkgmanager.Manager{{Name: "brew"}}
	plans := BuildFixPlans(report, managers)
	assert.Len(t, plans, 2, "should only plan fixes for docker and python")
}

func TestExecute_NoCommand(t *testing.T) {
	plan := FixPlan{
		Tool:       "unknown",
		ManualHint: "install it yourself",
	}
	result := Execute(plan)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "no command available")
}

func TestExecuteAll_EmptyPlans(t *testing.T) {
	results, ok := ExecuteAll(nil, false)
	assert.Empty(t, results)
	assert.True(t, ok)
}

func TestExecuteAll_NoCommandStopsOnError(t *testing.T) {
	plans := []FixPlan{
		{Tool: "a", ManualHint: "manual"},
		{Tool: "b", Command: []string{"echo", "hello"}, Manager: "test"},
	}
	results, ok := ExecuteAll(plans, false)
	assert.False(t, ok)
	assert.Len(t, results, 1, "should stop after first failure")
}

func TestExecuteAll_ContinueOnError(t *testing.T) {
	plans := []FixPlan{
		{Tool: "a", ManualHint: "manual"},
		{Tool: "b", ManualHint: "manual"},
	}
	results, ok := ExecuteAll(plans, true)
	assert.False(t, ok)
	assert.Len(t, results, 2, "should continue past failures")
}
