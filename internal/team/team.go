package team

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/groundctl/groundctl/internal/config"
	"github.com/groundctl/groundctl/internal/model"
)

// Pull fetches a team .ground.yaml from a source.
// Supports:
//   - Git repo URL: clones and extracts .ground.yaml
//   - Direct file URL: downloads the raw YAML
//   - Local path: copies from filesystem
func Pull(source, destPath string) (*model.GroundConfig, error) {
	if isLocalPath(source) {
		return pullLocal(source, destPath)
	}
	if isDirectURL(source) {
		return pullURL(source, destPath)
	}
	return pullGit(source, destPath)
}

// Push commits and pushes .ground.yaml to a git remote.
// Assumes the current directory is a git repo.
func Push(configPath, message string) error {
	if message == "" {
		message = "Update .ground.yaml via groundctl"
	}

	// Stage the file
	if err := gitRun("add", configPath); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// Commit
	if err := gitRun("commit", "-m", message); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	// Push
	if err := gitRun("push"); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}

func pullLocal(source, destPath string) (*model.GroundConfig, error) {
	cfg, err := config.Load(source)
	if err != nil {
		return nil, fmt.Errorf("could not load local config: %w", err)
	}
	if err := config.Save(destPath, cfg); err != nil {
		return nil, fmt.Errorf("could not save config: %w", err)
	}
	return cfg, nil
}

func pullURL(url, destPath string) (*model.GroundConfig, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: HTTP %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return nil, fmt.Errorf("could not write config: %w", err)
	}

	return config.Load(destPath)
}

func pullGit(repoURL, destPath string) (*model.GroundConfig, error) {
	// Clone into temp directory
	tmpDir, err := os.MkdirTemp("", "groundctl-pull-*")
	if err != nil {
		return nil, fmt.Errorf("could not create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Normalize URL
	repoURL = normalizeGitURL(repoURL)

	// Shallow clone
	cmd := exec.Command("git", "clone", "--depth=1", repoURL, tmpDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git clone failed: %s: %w", strings.TrimSpace(string(out)), err)
	}

	// Find .ground.yaml in the cloned repo
	srcPath := filepath.Join(tmpDir, config.DefaultConfigFile)
	if !config.Exists(srcPath) {
		return nil, fmt.Errorf("no %s found in repository %s", config.DefaultConfigFile, repoURL)
	}

	return pullLocal(srcPath, destPath)
}

func normalizeGitURL(url string) string {
	// Handle github.com/org/repo shorthand
	if !strings.Contains(url, "://") && !strings.HasPrefix(url, "git@") {
		if strings.Count(url, "/") >= 1 && !strings.HasPrefix(url, "/") {
			url = "https://" + url
		}
	}
	// Strip trailing .git if not present
	url = strings.TrimSuffix(url, "/")
	return url
}

func isLocalPath(source string) bool {
	if strings.HasPrefix(source, "/") || strings.HasPrefix(source, "./") || strings.HasPrefix(source, "..") {
		return true
	}
	// Windows paths
	if len(source) >= 2 && source[1] == ':' {
		return true
	}
	// Check if it's an existing file
	_, err := os.Stat(source)
	return err == nil
}

func isDirectURL(source string) bool {
	lower := strings.ToLower(source)
	return (strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")) &&
		(strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml"))
}

func gitRun(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
