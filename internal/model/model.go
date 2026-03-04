package model

import "time"

// Severity indicates whether a tool is required or recommended.
type Severity string

const (
	SeverityRequired    Severity = "required"
	SeverityRecommended Severity = "recommended"
)

// ToolSpec defines a tool requirement in .ground.yaml.
type ToolSpec struct {
	Name       string   `yaml:"name" json:"name"`
	Version    string   `yaml:"version,omitempty" json:"version,omitempty"`       // semver constraint e.g. ">=18.0.0", "^3.10"
	Severity   Severity `yaml:"severity,omitempty" json:"severity,omitempty"`     // required | recommended
	InstallCmd string   `yaml:"install_cmd,omitempty" json:"install_cmd,omitempty"` // manual install hint
}

// DetectedTool represents a tool found on the machine.
type DetectedTool struct {
	Name    string `json:"name"`
	Version string `json:"version"`           // raw version string detected
	Path    string `json:"path,omitempty"`     // binary path
	Found   bool   `json:"found"`
	Error   string `json:"error,omitempty"`    // detection error message
}

// Snapshot captures the full machine state at a point in time.
type Snapshot struct {
	Timestamp string         `json:"timestamp"`
	Hostname  string         `json:"hostname"`
	OS        string         `json:"os"`
	Arch      string         `json:"arch"`
	Tools     []DetectedTool `json:"tools"`
}

// DriftStatus represents the result of comparing a tool against its spec.
type DriftStatus string

const (
	DriftOK      DriftStatus = "ok"       // version matches constraint
	DriftWarning DriftStatus = "warning"  // recommended tool drifted or version mismatch
	DriftError   DriftStatus = "error"    // required tool missing or wrong version
)

// DriftItem represents the drift status of a single tool.
type DriftItem struct {
	Tool     string      `json:"tool"`
	Status   DriftStatus `json:"status"`
	Expected string      `json:"expected,omitempty"`  // version constraint from spec
	Actual   string      `json:"actual,omitempty"`    // detected version
	Message  string      `json:"message,omitempty"`
}

// DriftReport is the result of comparing a snapshot against a ground config.
type DriftReport struct {
	Timestamp string      `json:"timestamp"`
	Items     []DriftItem `json:"items"`
	Summary   Summary     `json:"summary"`
}

// Summary holds aggregate drift counts.
type Summary struct {
	Total    int `json:"total"`
	OK       int `json:"ok"`
	Warnings int `json:"warnings"`
	Errors   int `json:"errors"`
}

// SecretSpec defines a secret reference in .ground.yaml.
type SecretSpec struct {
	Name        string `yaml:"name" json:"name"`                                   // env var name, e.g. "DATABASE_URL"
	Ref         string `yaml:"ref" json:"ref"`                                     // secret reference, e.g. "${op://vault/item/field}"
	Description string `yaml:"description,omitempty" json:"description,omitempty"` // human-readable description
}

// GroundConfig represents the .ground.yaml file.
type GroundConfig struct {
	Name        string       `yaml:"name,omitempty" json:"name,omitempty"`
	Description string       `yaml:"description,omitempty" json:"description,omitempty"`
	Extends     string       `yaml:"extends,omitempty" json:"extends,omitempty"` // parent profile name for inheritance
	Team        *TeamMeta    `yaml:"team,omitempty" json:"team,omitempty"`
	Tools       []ToolSpec   `yaml:"tools" json:"tools"`
	Secrets     []SecretSpec `yaml:"secrets,omitempty" json:"secrets,omitempty"`
}

// TeamMeta holds team-level metadata for shared configs.
type TeamMeta struct {
	Org    string `yaml:"org,omitempty" json:"org,omitempty"`
	Repo   string `yaml:"repo,omitempty" json:"repo,omitempty"`
	Branch string `yaml:"branch,omitempty" json:"branch,omitempty"`
}

// MergeTools merges tools from a parent config into this config.
// Child tool specs override parent specs with the same name.
func (c *GroundConfig) MergeTools(parent *GroundConfig) {
	if parent == nil {
		return
	}
	existing := make(map[string]bool, len(c.Tools))
	for _, t := range c.Tools {
		existing[t.Name] = true
	}
	for _, t := range parent.Tools {
		if !existing[t.Name] {
			c.Tools = append(c.Tools, t)
		}
	}
}

// NewSnapshot creates a snapshot with the current timestamp.
func NewSnapshot(hostname, os, arch string, tools []DetectedTool) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Hostname:  hostname,
		OS:        os,
		Arch:      arch,
		Tools:     tools,
	}
}
