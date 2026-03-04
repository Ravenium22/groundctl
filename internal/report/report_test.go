package report

import (
	"testing"

	"github.com/groundctl/groundctl/internal/model"
	"github.com/stretchr/testify/assert"
)

func makeTestReport() *model.DriftReport {
	return &model.DriftReport{
		Timestamp: "2024-01-15T10:00:00Z",
		Items: []model.DriftItem{
			{Tool: "node", Status: model.DriftOK, Expected: ">=18.0.0", Actual: "20.10.0"},
			{Tool: "python", Status: model.DriftWarning, Expected: ">=3.10", Actual: "3.9.1", Message: "version drift"},
			{Tool: "docker", Status: model.DriftError, Expected: ">=24.0.0", Actual: "", Message: "not found"},
		},
		Summary: model.Summary{Total: 3, OK: 1, Warnings: 1, Errors: 1},
	}
}

func TestFormatJSON(t *testing.T) {
	report := makeTestReport()
	out, err := FormatJSON(report)
	assert.NoError(t, err)
	assert.Contains(t, out, `"tool": "node"`)
	assert.Contains(t, out, `"status": "ok"`)
	assert.Contains(t, out, `"total": 3`)
}

func TestFormatMarkdown(t *testing.T) {
	report := makeTestReport()
	out := FormatMarkdown(report)
	assert.Contains(t, out, "# groundctl Drift Report")
	assert.Contains(t, out, "| node | ok |")
	assert.Contains(t, out, "| python | warning |")
	assert.Contains(t, out, "| docker | error |")
	assert.Contains(t, out, "| Total | 3 |")
	assert.Contains(t, out, "| Warnings | 1 |")
}

func TestFormatHTML(t *testing.T) {
	report := makeTestReport()
	out := FormatHTML(report)
	assert.Contains(t, out, "<!DOCTYPE html>")
	assert.Contains(t, out, "groundctl Drift Report")
	assert.Contains(t, out, "status-ok")
	assert.Contains(t, out, "status-warning")
	assert.Contains(t, out, "status-error")
	assert.Contains(t, out, "<strong>node</strong>")
	assert.Contains(t, out, "<strong>docker</strong>")
}

func TestFormatMarkdownAllOK(t *testing.T) {
	report := &model.DriftReport{
		Items: []model.DriftItem{
			{Tool: "git", Status: model.DriftOK, Expected: ">=2.0.0", Actual: "2.43.0"},
		},
		Summary: model.Summary{Total: 1, OK: 1},
	}
	out := FormatMarkdown(report)
	assert.Contains(t, out, "| OK | 1 |")
	assert.Contains(t, out, "| Errors | 0 |")
}

func TestOrDash(t *testing.T) {
	assert.Equal(t, "-", orDash(""))
	assert.Equal(t, "hello", orDash("hello"))
}
