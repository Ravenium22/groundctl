package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/Ravenium22/groundctl/internal/cache"
	"github.com/Ravenium22/groundctl/internal/config"
	"github.com/Ravenium22/groundctl/internal/detector"
	"github.com/Ravenium22/groundctl/internal/drift"
	"github.com/Ravenium22/groundctl/internal/model"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check your machine against the team standard",
	Long: `Compares your machine's installed tools against the .ground.yaml standard
and shows a terraform-plan-style drift report.

Exit codes:
  0  All tools match the standard
  1  Warnings only (recommended tools drifted)
  2  Errors present (required tools missing or wrong version)`,
	RunE: runCheck,
}

var (
	checkJSON    bool
	checkQuiet   bool
	checkCI      bool
	checkNoCache bool
	checkConfig  string
)

func init() {
	checkCmd.Flags().BoolVar(&checkJSON, "json", false, "Output report as JSON")
	checkCmd.Flags().BoolVarP(&checkQuiet, "quiet", "q", false, "Only output errors and warnings")
	checkCmd.Flags().BoolVar(&checkCI, "ci", false, "Output GitHub Actions annotations")
	checkCmd.Flags().BoolVar(&checkNoCache, "no-cache", false, "Skip detection cache")
	checkCmd.Flags().StringVarP(&checkConfig, "config", "c", "", "Path to .ground.yaml (default: auto-detect)")
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	// Find config
	configPath := checkConfig
	if configPath == "" {
		configPath = findConfig()
	}
	if configPath == "" {
		return fmt.Errorf("no .ground.yaml found. Run 'ground init' to create one")
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	if len(cfg.Tools) == 0 {
		return fmt.Errorf("no tools defined in %s", configPath)
	}

	// Detect only the tools listed in the config
	names := make([]string, len(cfg.Tools))
	for i, t := range cfg.Tools {
		names[i] = t.Name
	}

	var detected []model.DetectedTool
	if checkNoCache {
		detected = detector.DetectByNames(names)
	} else {
		detected = detectWithCache(names)
	}

	// Run diff engine
	report := drift.Compare(cfg, detected)

	// Output
	ciMode := checkCI || os.Getenv("GITHUB_ACTIONS") == "true"
	if checkJSON {
		return outputJSON(report)
	}
	if ciMode {
		outputCIAnnotations(report, configPath)
	} else {
		outputPretty(report, configPath, checkQuiet)
	}

	// Exit with appropriate code
	code := drift.ExitCode(report)
	if code != 0 {
		os.Exit(code)
	}
	return nil
}

func findConfig() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Walk up directories looking for .ground.yaml
	for {
		path := filepath.Join(dir, config.DefaultConfigFile)
		if config.Exists(path) {
			return path
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func outputJSON(report *model.DriftReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func outputPretty(report *model.DriftReport, configPath string, quiet bool) {
	// Styles
	titleStyle := lipgloss.NewStyle().Bold(true)
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	nameStyle := lipgloss.NewStyle().Bold(true).Width(20)

	// Header
	fmt.Println()
	fmt.Printf("%s %s\n", titleStyle.Render("groundctl"), dimStyle.Render("drift report"))
	fmt.Printf("%s %s\n", dimStyle.Render("config:"), configPath)
	fmt.Println()

	// Items
	for _, item := range report.Items {
		if quiet && item.Status == model.DriftOK {
			continue
		}

		var icon string
		var style lipgloss.Style
		switch item.Status {
		case model.DriftOK:
			icon = okStyle.Render("[ok]")
			style = okStyle
		case model.DriftWarning:
			icon = warnStyle.Render("[!!]")
			style = warnStyle
		case model.DriftError:
			icon = errStyle.Render("[ERR]")
			style = errStyle
		}

		name := nameStyle.Render(item.Tool)

		var detail string
		switch {
		case item.Status == model.DriftOK && item.Actual != "":
			detail = dimStyle.Render(item.Actual)
		case item.Message != "":
			detail = style.Render(item.Message)
			if item.Actual != "" {
				detail += dimStyle.Render(fmt.Sprintf(" (have: %s, want: %s)", item.Actual, item.Expected))
			}
		}

		fmt.Printf("  %s %s %s\n", icon, name, detail)
	}

	// Summary
	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))

	parts := []string{
		fmt.Sprintf("%d checked", report.Summary.Total),
	}
	if report.Summary.OK > 0 {
		parts = append(parts, okStyle.Render(fmt.Sprintf("%d ok", report.Summary.OK)))
	}
	if report.Summary.Warnings > 0 {
		parts = append(parts, warnStyle.Render(fmt.Sprintf("%d warnings", report.Summary.Warnings)))
	}
	if report.Summary.Errors > 0 {
		parts = append(parts, errStyle.Render(fmt.Sprintf("%d errors", report.Summary.Errors)))
	}

	fmt.Printf("  %s\n", strings.Join(parts, "  "))

	if report.Summary.Errors > 0 {
		fmt.Println()
		fmt.Printf("  %s\n", errStyle.Render("Run 'ground fix' to resolve drift."))
	} else if report.Summary.Warnings > 0 {
		fmt.Println()
		fmt.Printf("  %s\n", warnStyle.Render("Some recommended tools have drifted."))
	} else {
		fmt.Println()
		fmt.Printf("  %s\n", okStyle.Render("All tools match the team standard."))
	}
	fmt.Println()
}

func detectWithCache(names []string) []model.DetectedTool {
	store := cache.New()
	_ = store.Load()

	var uncached []string
	cachedResults := make(map[string]model.DetectedTool)

	for _, name := range names {
		if tool, ok := store.Get(name); ok {
			cachedResults[name] = tool
		} else {
			uncached = append(uncached, name)
		}
	}

	// Detect only uncached tools
	var freshResults []model.DetectedTool
	if len(uncached) > 0 {
		freshResults = detector.DetectByNames(uncached)
		for _, tool := range freshResults {
			store.Put(tool)
		}
	}

	_ = store.Save()

	// Merge results in original order
	freshMap := make(map[string]model.DetectedTool, len(freshResults))
	for _, t := range freshResults {
		freshMap[t.Name] = t
	}

	results := make([]model.DetectedTool, 0, len(names))
	for _, name := range names {
		if t, ok := cachedResults[name]; ok {
			results = append(results, t)
		} else if t, ok := freshMap[name]; ok {
			results = append(results, t)
		}
	}
	return results
}

func outputCIAnnotations(report *model.DriftReport, configPath string) {
	for _, item := range report.Items {
		switch item.Status {
		case model.DriftError:
			msg := item.Message
			if item.Actual != "" {
				msg += fmt.Sprintf(" (have: %s, want: %s)", item.Actual, item.Expected)
			}
			fmt.Printf("::error file=%s,title=groundctl: %s::%s\n", configPath, item.Tool, msg)
		case model.DriftWarning:
			msg := item.Message
			if item.Actual != "" {
				msg += fmt.Sprintf(" (have: %s, want: %s)", item.Actual, item.Expected)
			}
			fmt.Printf("::warning file=%s,title=groundctl: %s::%s\n", configPath, item.Tool, msg)
		case model.DriftOK:
			fmt.Printf("::notice file=%s,title=groundctl: %s::ok (%s)\n", configPath, item.Tool, item.Actual)
		}
	}

	// Summary as group
	fmt.Printf("::group::groundctl drift summary\n")
	fmt.Printf("%d checked, %d ok, %d warnings, %d errors\n",
		report.Summary.Total, report.Summary.OK, report.Summary.Warnings, report.Summary.Errors)
	fmt.Printf("::endgroup::\n")
}
