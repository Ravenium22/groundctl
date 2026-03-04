package fixer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/groundctl/groundctl/internal/drift"
	"github.com/groundctl/groundctl/internal/model"
	"github.com/groundctl/groundctl/internal/pkgmanager"
)

// FixResult captures the outcome of a single fix attempt.
type FixResult struct {
	Plan    FixPlan `json:"plan"`
	Success bool    `json:"success"`
	Output  string  `json:"output,omitempty"`
	Error   string  `json:"error,omitempty"`
}

// BuildFixPlans takes a drift report and available PMs, producing fix plans
// for all items that need fixing (errors and warnings).
func BuildFixPlans(report *model.DriftReport, managers []pkgmanager.Manager) []FixPlan {
	pmNames := make([]string, len(managers))
	for i, m := range managers {
		pmNames[i] = m.Name
	}

	var plans []FixPlan
	for _, item := range report.Items {
		if item.Status == model.DriftOK {
			continue
		}

		action := ActionInstall
		if item.Actual != "" {
			action = ActionUpgrade
		}

		plan := ResolveStrategy(item.Tool, action, pmNames, item.Expected, item.Actual)
		plans = append(plans, plan)
	}
	return plans
}

// Execute runs a single fix plan and returns the result.
func Execute(plan FixPlan) FixResult {
	if len(plan.Command) == 0 {
		return FixResult{
			Plan:    plan,
			Success: false,
			Error:   "no command available; " + plan.ManualHint,
		}
	}

	cmd := exec.Command(plan.Command[0], plan.Command[1:]...)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	if err != nil {
		return FixResult{
			Plan:    plan,
			Success: false,
			Output:  output,
			Error:   fmt.Sprintf("command failed: %v", err),
		}
	}

	return FixResult{
		Plan:    plan,
		Success: true,
		Output:  output,
	}
}

// ExecuteAll runs all fix plans sequentially, stopping on first failure
// unless continueOnError is true. Returns results and overall success.
func ExecuteAll(plans []FixPlan, continueOnError bool) ([]FixResult, bool) {
	var results []FixResult
	allSuccess := true

	for _, plan := range plans {
		if len(plan.Command) == 0 {
			results = append(results, FixResult{
				Plan:    plan,
				Success: false,
				Error:   "no command available; " + plan.ManualHint,
			})
			allSuccess = false
			if !continueOnError {
				break
			}
			continue
		}

		result := Execute(plan)
		results = append(results, result)

		if !result.Success {
			allSuccess = false
			if !continueOnError {
				break
			}
		}
	}

	return results, allSuccess
}

// Verify runs a fresh drift check after fixes to confirm resolution.
func Verify(cfg *model.GroundConfig, detected []model.DetectedTool) *model.DriftReport {
	return drift.Compare(cfg, detected)
}
