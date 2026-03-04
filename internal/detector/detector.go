package detector

import (
	"context"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Ravenium22/groundctl/internal/model"
)

// ToolDetector knows how to detect a specific tool.
type ToolDetector struct {
	Name      string
	Commands  [][]string // commands to try, in order
	VersionRe *regexp.Regexp
}

// DetectTimeout is the maximum time to wait for a single tool detection.
var DetectTimeout = 5 * time.Second

// registry holds all known tool detectors.
var registry []ToolDetector

func init() {
	semverRe := regexp.MustCompile(`(\d+\.\d+(?:\.\d+)*)`)

	registry = []ToolDetector{
		{Name: "node", Commands: [][]string{{"node", "--version"}}, VersionRe: semverRe},
		{Name: "npm", Commands: [][]string{{"npm", "--version"}}, VersionRe: semverRe},
		{Name: "python", Commands: pythonCommands(), VersionRe: semverRe},
		{Name: "pip", Commands: [][]string{{"pip", "--version"}, {"pip3", "--version"}}, VersionRe: semverRe},
		{Name: "go", Commands: [][]string{{"go", "version"}}, VersionRe: semverRe},
		{Name: "git", Commands: [][]string{{"git", "--version"}}, VersionRe: semverRe},
		{Name: "docker", Commands: [][]string{{"docker", "--version"}}, VersionRe: semverRe},
		{Name: "docker-compose", Commands: [][]string{{"docker", "compose", "version"}, {"docker-compose", "--version"}}, VersionRe: semverRe},
		{Name: "kubectl", Commands: [][]string{{"kubectl", "version", "--client", "--short"}}, VersionRe: semverRe},
		{Name: "terraform", Commands: [][]string{{"terraform", "--version"}}, VersionRe: semverRe},
		{Name: "java", Commands: [][]string{{"java", "-version"}}, VersionRe: semverRe},
		{Name: "ruby", Commands: [][]string{{"ruby", "--version"}}, VersionRe: semverRe},
		{Name: "rustc", Commands: [][]string{{"rustc", "--version"}}, VersionRe: semverRe},
		{Name: "cargo", Commands: [][]string{{"cargo", "--version"}}, VersionRe: semverRe},
		{Name: "make", Commands: [][]string{{"make", "--version"}}, VersionRe: semverRe},
		{Name: "gh", Commands: [][]string{{"gh", "--version"}}, VersionRe: semverRe},
		{Name: "curl", Commands: [][]string{{"curl", "--version"}}, VersionRe: semverRe},
		{Name: "wget", Commands: [][]string{{"wget", "--version"}}, VersionRe: semverRe},
	}
}

func pythonCommands() [][]string {
	if runtime.GOOS == "windows" {
		return [][]string{{"python", "--version"}, {"python3", "--version"}, {"py", "--version"}}
	}
	return [][]string{{"python3", "--version"}, {"python", "--version"}}
}

// DetectAll runs all registered detectors concurrently and returns results.
func DetectAll() []model.DetectedTool {
	return detectParallel(registry)
}

// DetectByNames runs only the detectors matching the given tool names, concurrently.
func DetectByNames(names []string) []model.DetectedTool {
	nameSet := make(map[string]bool, len(names))
	for _, n := range names {
		nameSet[strings.ToLower(n)] = true
	}

	var filtered []ToolDetector
	for _, d := range registry {
		if nameSet[d.Name] {
			filtered = append(filtered, d)
		}
	}

	return detectParallel(filtered)
}

// detectParallel runs detections concurrently with bounded parallelism.
func detectParallel(detectors []ToolDetector) []model.DetectedTool {
	if len(detectors) == 0 {
		return nil
	}

	results := make([]model.DetectedTool, len(detectors))
	var wg sync.WaitGroup

	// Limit concurrency to NumCPU
	maxWorkers := runtime.NumCPU()
	if maxWorkers > len(detectors) {
		maxWorkers = len(detectors)
	}
	sem := make(chan struct{}, maxWorkers)

	for i, d := range detectors {
		wg.Add(1)
		go func(idx int, det ToolDetector) {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release
			results[idx] = detect(det)
		}(i, d)
	}

	wg.Wait()
	return results
}

// ListKnownTools returns the names of all tools we can detect.
func ListKnownTools() []string {
	names := make([]string, len(registry))
	for i, d := range registry {
		names[i] = d.Name
	}
	return names
}

func detect(d ToolDetector) model.DetectedTool {
	for _, args := range d.Commands {
		ctx, cancel := context.WithTimeout(context.Background(), DetectTimeout)
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		out, err := cmd.CombinedOutput()
		cancel()

		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return model.DetectedTool{
					Name:  d.Name,
					Found: false,
					Error: "detection timed out",
				}
			}
			continue
		}

		output := strings.TrimSpace(string(out))
		if match := d.VersionRe.FindString(output); match != "" {
			path, _ := exec.LookPath(args[0])
			return model.DetectedTool{
				Name:    d.Name,
				Version: match,
				Path:    path,
				Found:   true,
			}
		}
	}

	return model.DetectedTool{
		Name:  d.Name,
		Found: false,
		Error: "not found or version could not be detected",
	}
}

// ParseVersion extracts a version string from raw command output.
func ParseVersion(output string) string {
	re := regexp.MustCompile(`(\d+\.\d+(?:\.\d+)*)`)
	if match := re.FindString(output); match != "" {
		return match
	}
	return ""
}
