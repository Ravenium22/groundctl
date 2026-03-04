package drift

import (
	"fmt"
	"time"

	"github.com/groundctl/groundctl/internal/model"
	"github.com/groundctl/groundctl/internal/version"
)

// Compare takes a ground config and a list of detected tools,
// producing a DriftReport showing what matches, what drifted, and what's missing.
func Compare(cfg *model.GroundConfig, detected []model.DetectedTool) *model.DriftReport {
	// Index detected tools by name for O(1) lookup
	detectedMap := make(map[string]model.DetectedTool, len(detected))
	for _, d := range detected {
		detectedMap[d.Name] = d
	}

	report := &model.DriftReport{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	for _, spec := range cfg.Tools {
		item := compareOne(spec, detectedMap)
		report.Items = append(report.Items, item)

		report.Summary.Total++
		switch item.Status {
		case model.DriftOK:
			report.Summary.OK++
		case model.DriftWarning:
			report.Summary.Warnings++
		case model.DriftError:
			report.Summary.Errors++
		}
	}

	return report
}

func compareOne(spec model.ToolSpec, detected map[string]model.DetectedTool) model.DriftItem {
	tool, found := detected[spec.Name]

	severity := spec.Severity
	if severity == "" {
		severity = model.SeverityRequired
	}

	// Tool not found on machine
	if !found || !tool.Found {
		status := model.DriftError
		if severity == model.SeverityRecommended {
			status = model.DriftWarning
		}
		return model.DriftItem{
			Tool:     spec.Name,
			Status:   status,
			Expected: spec.Version,
			Message:  "not installed",
		}
	}

	// No version constraint specified - just check presence
	if spec.Version == "" || spec.Version == "*" {
		return model.DriftItem{
			Tool:   spec.Name,
			Status: model.DriftOK,
			Actual: tool.Version,
		}
	}

	// Check version constraint
	satisfied, err := version.CheckConstraint(spec.Version, tool.Version)
	if err != nil {
		return model.DriftItem{
			Tool:     spec.Name,
			Status:   model.DriftWarning,
			Expected: spec.Version,
			Actual:   tool.Version,
			Message:  fmt.Sprintf("could not parse version: %v", err),
		}
	}

	if satisfied {
		return model.DriftItem{
			Tool:     spec.Name,
			Status:   model.DriftOK,
			Expected: spec.Version,
			Actual:   tool.Version,
		}
	}

	// Version doesn't satisfy constraint
	status := model.DriftError
	if severity == model.SeverityRecommended {
		status = model.DriftWarning
	}
	return model.DriftItem{
		Tool:     spec.Name,
		Status:   status,
		Expected: spec.Version,
		Actual:   tool.Version,
		Message:  fmt.Sprintf("version %s does not satisfy %s", tool.Version, spec.Version),
	}
}

// ExitCode returns the appropriate process exit code for a drift report.
//   - 0: all ok
//   - 1: warnings only
//   - 2: errors present
func ExitCode(report *model.DriftReport) int {
	if report.Summary.Errors > 0 {
		return 2
	}
	if report.Summary.Warnings > 0 {
		return 1
	}
	return 0
}
