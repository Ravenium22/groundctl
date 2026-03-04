package config

import (
	"fmt"
	"strings"

	"github.com/Ravenium22/groundctl/internal/model"
	"github.com/Ravenium22/groundctl/internal/version"
)

// ValidationError represents a single validation issue.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult holds all validation errors found.
type ValidationResult struct {
	Errors []ValidationError `json:"errors"`
}

// IsValid returns true if no validation errors were found.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// Add appends a validation error.
func (r *ValidationResult) Add(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

// Error returns all errors as a single string.
func (r *ValidationResult) Error() string {
	if r.IsValid() {
		return ""
	}
	var parts []string
	for _, e := range r.Errors {
		parts = append(parts, e.Error())
	}
	return strings.Join(parts, "; ")
}

// Validate checks a GroundConfig for common issues.
func Validate(cfg *model.GroundConfig) *ValidationResult {
	result := &ValidationResult{}

	if len(cfg.Tools) == 0 {
		result.Add("tools", "at least one tool must be defined")
	}

	seen := make(map[string]bool)
	for i, tool := range cfg.Tools {
		field := fmt.Sprintf("tools[%d]", i)

		// Name required
		if tool.Name == "" {
			result.Add(field+".name", "tool name is required")
			continue
		}

		// Duplicate check
		if seen[tool.Name] {
			result.Add(field+".name", fmt.Sprintf("duplicate tool %q", tool.Name))
		}
		seen[tool.Name] = true

		// Severity must be valid
		if tool.Severity != "" &&
			tool.Severity != model.SeverityRequired &&
			tool.Severity != model.SeverityRecommended {
			result.Add(field+".severity",
				fmt.Sprintf("invalid severity %q (must be %q or %q)",
					tool.Severity, model.SeverityRequired, model.SeverityRecommended))
		}

		// Version constraint must be parseable
		if tool.Version != "" && tool.Version != "*" {
			if _, err := version.CheckConstraint(tool.Version, "99.99.99"); err != nil {
				result.Add(field+".version",
					fmt.Sprintf("invalid version constraint %q: %v", tool.Version, err))
			}
		}
	}

	// Validate secrets
	secretSeen := make(map[string]bool)
	for i, s := range cfg.Secrets {
		field := fmt.Sprintf("secrets[%d]", i)

		if s.Name == "" {
			result.Add(field+".name", "secret name is required")
			continue
		}

		if secretSeen[s.Name] {
			result.Add(field+".name", fmt.Sprintf("duplicate secret %q", s.Name))
		}
		secretSeen[s.Name] = true

		if s.Ref == "" {
			result.Add(field+".ref", "secret reference is required")
		} else if !strings.Contains(s.Ref, "://") || !strings.HasPrefix(s.Ref, "${") {
			result.Add(field+".ref",
				fmt.Sprintf("invalid secret reference %q (expected ${backend://path} format)", s.Ref))
		}
	}

	return result
}
