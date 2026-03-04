package fixer

import "runtime"

// Action describes what to do to fix a tool.
type Action string

const (
	ActionInstall Action = "install"
	ActionUpgrade Action = "upgrade"
)

// FixPlan describes a single fix operation.
type FixPlan struct {
	Tool       string   `json:"tool"`
	Action     Action   `json:"action"`
	Manager    string   `json:"manager"`              // package manager to use
	Command    []string `json:"command"`               // full command to execute
	CommandStr string   `json:"command_str"`           // human-readable command string
	Expected   string   `json:"expected,omitempty"`    // desired version constraint
	Current    string   `json:"current,omitempty"`     // current version (empty if missing)
	ManualHint string   `json:"manual_hint,omitempty"` // fallback instructions if no PM
}

// installMap defines how to install each tool per package manager.
// Key: tool name, Value: map of PM name -> package name/command args.
var installMap = map[string]map[string]string{
	"node":           {"brew": "node", "apt": "nodejs", "winget": "OpenJS.NodeJS.LTS", "scoop": "nodejs-lts", "choco": "nodejs-lts", "dnf": "nodejs", "pacman": "nodejs"},
	"npm":            {"brew": "node", "apt": "npm", "winget": "OpenJS.NodeJS.LTS", "scoop": "nodejs-lts", "choco": "nodejs-lts", "dnf": "npm", "pacman": "npm"},
	"python":         {"brew": "python", "apt": "python3", "winget": "Python.Python.3.12", "scoop": "python", "choco": "python3", "dnf": "python3", "pacman": "python"},
	"pip":            {"brew": "python", "apt": "python3-pip", "winget": "Python.Python.3.12", "scoop": "python", "choco": "python3", "dnf": "python3-pip", "pacman": "python-pip"},
	"go":             {"brew": "go", "apt": "golang-go", "winget": "GoLang.Go", "scoop": "go", "choco": "golang", "dnf": "golang", "pacman": "go"},
	"git":            {"brew": "git", "apt": "git", "winget": "Git.Git", "scoop": "git", "choco": "git", "dnf": "git", "pacman": "git"},
	"docker":         {"brew": "docker", "apt": "docker.io", "winget": "Docker.DockerDesktop", "scoop": "docker", "choco": "docker-desktop", "dnf": "docker-ce", "pacman": "docker"},
	"docker-compose": {"brew": "docker-compose", "apt": "docker-compose-v2", "winget": "Docker.DockerDesktop", "choco": "docker-compose", "dnf": "docker-compose-plugin", "pacman": "docker-compose"},
	"kubectl":        {"brew": "kubectl", "apt": "kubectl", "winget": "Kubernetes.kubectl", "scoop": "kubectl", "choco": "kubernetes-cli", "dnf": "kubectl", "pacman": "kubectl"},
	"terraform":      {"brew": "terraform", "apt": "terraform", "winget": "Hashicorp.Terraform", "scoop": "terraform", "choco": "terraform", "dnf": "terraform", "pacman": "terraform"},
	"java":           {"brew": "openjdk", "apt": "default-jdk", "winget": "EclipseAdoptium.Temurin.21.JDK", "scoop": "temurin-jdk", "choco": "temurin", "dnf": "java-latest-openjdk", "pacman": "jdk-openjdk"},
	"ruby":           {"brew": "ruby", "apt": "ruby-full", "winget": "RubyInstallerTeam.Ruby.3.3", "scoop": "ruby", "choco": "ruby", "dnf": "ruby", "pacman": "ruby"},
	"rustc":          {"brew": "rust", "apt": "rustc", "winget": "Rustlang.Rustup", "scoop": "rustup", "choco": "rustup.install", "dnf": "rust", "pacman": "rust"},
	"cargo":          {"brew": "rust", "apt": "cargo", "winget": "Rustlang.Rustup", "scoop": "rustup", "choco": "rustup.install", "dnf": "cargo", "pacman": "rust"},
	"make":           {"brew": "make", "apt": "make", "winget": "GnuWin32.Make", "scoop": "make", "choco": "make", "dnf": "make", "pacman": "make"},
	"gh":             {"brew": "gh", "apt": "gh", "winget": "GitHub.cli", "scoop": "gh", "choco": "gh", "dnf": "gh", "pacman": "github-cli"},
	"curl":           {"brew": "curl", "apt": "curl", "winget": "cURL.cURL", "scoop": "curl", "choco": "curl", "dnf": "curl", "pacman": "curl"},
	"wget":           {"brew": "wget", "apt": "wget", "winget": "JernejSimoncic.Wget", "scoop": "wget", "choco": "wget", "dnf": "wget", "pacman": "wget"},
}

// installCmdTemplate defines the install command format per package manager.
var installCmdTemplate = map[string][]string{
	"brew":   {"brew", "install"},
	"apt":    {"sudo", "apt-get", "install", "-y"},
	"winget": {"winget", "install", "--accept-source-agreements", "--accept-package-agreements"},
	"scoop":  {"scoop", "install"},
	"choco":  {"choco", "install", "-y"},
	"dnf":    {"sudo", "dnf", "install", "-y"},
	"pacman": {"sudo", "pacman", "-S", "--noconfirm"},
}

// upgradeCmdTemplate defines the upgrade command format per package manager.
var upgradeCmdTemplate = map[string][]string{
	"brew":   {"brew", "upgrade"},
	"apt":    {"sudo", "apt-get", "install", "--only-upgrade", "-y"},
	"winget": {"winget", "upgrade", "--accept-source-agreements", "--accept-package-agreements"},
	"scoop":  {"scoop", "update"},
	"choco":  {"choco", "upgrade", "-y"},
	"dnf":    {"sudo", "dnf", "upgrade", "-y"},
	"pacman": {"sudo", "pacman", "-S", "--noconfirm"},
}

// manualHints provides fallback install instructions.
var manualHints = map[string]string{
	"node":    "Visit https://nodejs.org/ or use nvm: curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash",
	"python":  "Visit https://python.org/downloads/",
	"go":      "Visit https://go.dev/dl/",
	"docker":  "Visit https://docs.docker.com/get-docker/",
	"rustc":   "Run: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh",
	"cargo":   "Run: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh",
	"kubectl": "Visit https://kubernetes.io/docs/tasks/tools/",
	"terraform": "Visit https://developer.hashicorp.com/terraform/downloads",
}

// ResolveStrategy determines the best fix plan for a tool given available package managers.
func ResolveStrategy(toolName string, action Action, availablePMs []string, expected, current string) FixPlan {
	return ResolveStrategyForPlatform(toolName, action, availablePMs, expected, current, runtime.GOOS)
}

// ResolveStrategyForPlatform is testable variant that accepts a platform string.
func ResolveStrategyForPlatform(toolName string, action Action, availablePMs []string, expected, current, platform string) FixPlan {
	plan := FixPlan{
		Tool:     toolName,
		Action:   action,
		Expected: expected,
		Current:  current,
	}

	pkgMap, hasTool := installMap[toolName]
	if !hasTool {
		plan.ManualHint = "No install strategy known for " + toolName
		return plan
	}

	// Try each available PM in priority order
	for _, pm := range availablePMs {
		pkg, ok := pkgMap[pm]
		if !ok {
			continue
		}

		var tmpl []string
		if action == ActionUpgrade {
			tmpl = upgradeCmdTemplate[pm]
		} else {
			tmpl = installCmdTemplate[pm]
		}
		if tmpl == nil {
			continue
		}

		cmd := make([]string, len(tmpl)+1)
		copy(cmd, tmpl)
		cmd[len(tmpl)] = pkg

		plan.Manager = pm
		plan.Command = cmd
		plan.CommandStr = formatCommand(cmd)
		return plan
	}

	// No PM found — provide manual hint
	if hint, ok := manualHints[toolName]; ok {
		plan.ManualHint = hint
	} else {
		plan.ManualHint = "Install " + toolName + " manually — no supported package manager found"
	}
	return plan
}

// LookupPackageName returns the package name for a tool in a given PM.
func LookupPackageName(toolName, pmName string) (string, bool) {
	pkgs, ok := installMap[toolName]
	if !ok {
		return "", false
	}
	pkg, ok := pkgs[pmName]
	return pkg, ok
}

func formatCommand(cmd []string) string {
	s := ""
	for i, part := range cmd {
		if i > 0 {
			s += " "
		}
		s += part
	}
	return s
}
