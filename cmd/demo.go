package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/Ravenium22/groundctl/internal/config"
	"github.com/Ravenium22/groundctl/internal/detector"
	"github.com/Ravenium22/groundctl/internal/drift"
	"github.com/Ravenium22/groundctl/internal/model"
	"github.com/spf13/cobra"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run an interactive demo of groundctl",
	Long: `Runs a guided walkthrough showing groundctl's core features:
scan -> check -> report. Takes about 30 seconds.`,
	RunE: runDemo,
}

func init() {
	rootCmd.AddCommand(demoCmd)
}

func runDemo(cmd *cobra.Command, args []string) error {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	nameStyle := lipgloss.NewStyle().Bold(true).Width(20)

	// Header
	fmt.Println()
	fmt.Printf("  %s\n", title.Render("groundctl demo"))
	fmt.Printf("  %s\n", dimStyle.Render("terraform plan for your local developer machine"))
	fmt.Println()

	// Step 1: Scan
	fmt.Printf("  %s Scanning your machine for installed tools...\n", dimStyle.Render("[1/3]"))
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	tools := detector.DetectAll()
	var found []model.DetectedTool
	for _, t := range tools {
		if t.Found {
			found = append(found, t)
		}
	}

	for _, t := range found {
		fmt.Printf("    %s %s %s\n", okStyle.Render("found"), nameStyle.Render(t.Name), dimStyle.Render(t.Version))
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf("  %s\n", okStyle.Render(fmt.Sprintf("Detected %d tools on your machine.", len(found))))
	fmt.Println()
	time.Sleep(500 * time.Millisecond)

	// Step 2: Create a demo config and check against it
	fmt.Printf("  %s Checking against a team standard...\n", dimStyle.Render("[2/3]"))
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// Build a demo config from detected tools with some artificial drift
	demoCfg := &model.GroundConfig{
		Name: "demo-team-standard",
		Tools: []model.ToolSpec{},
	}

	for i, t := range found {
		spec := model.ToolSpec{
			Name:     t.Name,
			Version:  fmt.Sprintf(">=%s", t.Version),
			Severity: model.SeverityRequired,
		}
		// Make one tool show as a warning (version "too new")
		if i == len(found)-1 && len(found) > 2 {
			spec.Version = ">=99.0.0"
			spec.Severity = model.SeverityRecommended
		}
		demoCfg.Tools = append(demoCfg.Tools, spec)
	}

	// Add a fake missing required tool
	demoCfg.Tools = append(demoCfg.Tools, model.ToolSpec{
		Name:     "helm",
		Version:  ">=3.0.0",
		Severity: model.SeverityRequired,
	})

	report := drift.Compare(demoCfg, append(found, model.DetectedTool{Name: "helm", Found: false, Error: "not found"}))

	for _, item := range report.Items {
		var icon string
		switch item.Status {
		case model.DriftOK:
			icon = okStyle.Render("[ok]")
		case model.DriftWarning:
			icon = warnStyle.Render("[!!]")
		case model.DriftError:
			icon = errStyle.Render("[ERR]")
		}

		detail := ""
		if item.Status == model.DriftOK {
			detail = dimStyle.Render(item.Actual)
		} else if item.Message != "" {
			detail = item.Message
			if item.Actual != "" {
				detail += dimStyle.Render(fmt.Sprintf(" (have: %s)", item.Actual))
			}
		}
		fmt.Printf("    %s %s %s\n", icon, nameStyle.Render(item.Tool), detail)
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf("  %s checked, %s, %s, %s\n",
		dimStyle.Render(fmt.Sprintf("%d", report.Summary.Total)),
		okStyle.Render(fmt.Sprintf("%d ok", report.Summary.OK)),
		warnStyle.Render(fmt.Sprintf("%d warnings", report.Summary.Warnings)),
		errStyle.Render(fmt.Sprintf("%d errors", report.Summary.Errors)))
	fmt.Println()
	time.Sleep(500 * time.Millisecond)

	// Step 3: What you'd do next
	fmt.Printf("  %s What comes next...\n", dimStyle.Render("[3/3]"))
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	steps := []struct{ cmd, desc string }{
		{"ground init", "Create .ground.yaml from your machine"},
		{"ground check", "Compare against the team standard"},
		{"ground fix", "Auto-fix all detected drift"},
		{"ground pull <repo>", "Fetch a team standard from git"},
		{"ground doctor", "Diagnose your setup"},
	}

	for _, s := range steps {
		fmt.Printf("    %s  %s\n",
			okStyle.Render(fmt.Sprintf("%-25s", s.cmd)),
			dimStyle.Render(s.desc))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf("  %s\n", title.Render("Ready to get started?"))
	fmt.Println()

	// Check if .ground.yaml exists
	if configPath := findConfig(); configPath != "" {
		fmt.Printf("  %s .ground.yaml found at %s\n", okStyle.Render("tip:"), configPath)
		fmt.Printf("  %s\n", dimStyle.Render("Run 'ground check' to see your drift report."))
	} else {
		fmt.Printf("  %s\n", dimStyle.Render("Run 'ground init' to create your first .ground.yaml"))
	}

	// Save demo config if user doesn't have one
	if findConfig() == "" {
		demoPath := filepath.Join(".", config.DefaultConfigFile)
		if err := config.Save(demoPath, demoCfg); err == nil {
			fmt.Printf("  %s\n", dimStyle.Render(fmt.Sprintf("Demo config saved to %s", demoPath)))
		}
	}

	fmt.Println()
	return nil
}
