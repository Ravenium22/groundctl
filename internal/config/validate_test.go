package config

import (
	"testing"

	"github.com/Ravenium22/groundctl/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestValidate_ValidConfig(t *testing.T) {
	cfg := &model.GroundConfig{
		Name: "test",
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=18.0.0", Severity: model.SeverityRequired},
			{Name: "go", Version: "^1.21", Severity: model.SeverityRecommended},
			{Name: "git"},
		},
	}
	result := Validate(cfg)
	assert.True(t, result.IsValid(), "expected valid, got: %s", result.Error())
}

func TestValidate_EmptyTools(t *testing.T) {
	cfg := &model.GroundConfig{Name: "test", Tools: nil}
	result := Validate(cfg)
	assert.False(t, result.IsValid())
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0].Message, "at least one tool")
}

func TestValidate_EmptyToolName(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "", Version: ">=1.0.0"},
		},
	}
	result := Validate(cfg)
	assert.False(t, result.IsValid())
	assert.Contains(t, result.Error(), "name is required")
}

func TestValidate_DuplicateTools(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">=18"},
			{Name: "node", Version: ">=20"},
		},
	}
	result := Validate(cfg)
	assert.False(t, result.IsValid())
	assert.Contains(t, result.Error(), "duplicate")
}

func TestValidate_InvalidSeverity(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "node", Severity: "critical"},
		},
	}
	result := Validate(cfg)
	assert.False(t, result.IsValid())
	assert.Contains(t, result.Error(), "invalid severity")
}

func TestValidate_InvalidVersion(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "node", Version: ">>>invalid"},
		},
	}
	result := Validate(cfg)
	assert.False(t, result.IsValid())
	assert.Contains(t, result.Error(), "invalid version constraint")
}

func TestValidate_StarVersionOK(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "git", Version: "*"},
		},
	}
	result := Validate(cfg)
	assert.True(t, result.IsValid())
}

func TestValidate_NoVersionOK(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: "git"},
		},
	}
	result := Validate(cfg)
	assert.True(t, result.IsValid())
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := &model.GroundConfig{
		Tools: []model.ToolSpec{
			{Name: ""},
			{Name: "node", Version: ">>>bad"},
			{Name: "go", Severity: "wrong"},
			{Name: "go"}, // duplicate
		},
	}
	result := Validate(cfg)
	assert.False(t, result.IsValid())
	assert.True(t, len(result.Errors) >= 3, "expected at least 3 errors, got %d", len(result.Errors))
}

func TestValidationResult_Error(t *testing.T) {
	r := &ValidationResult{}
	assert.Equal(t, "", r.Error())

	r.Add("field1", "msg1")
	r.Add("field2", "msg2")
	assert.Contains(t, r.Error(), "field1: msg1")
	assert.Contains(t, r.Error(), "field2: msg2")
}
