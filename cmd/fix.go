package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/Ravenium22/groundctl/internal/config"
	"github.com/Ravenium22/groundctl/internal/detector"
	"github.com/Ravenium22/groundctl/internal/drift"
	"github.com/Ravenium22/groundctl/internal/fixer"
	"github.com/Ravenium22/groundctl/internal/pkgmanager"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Auto-fix detected drift",
	Long: `Detects drift against .ground.yaml and resolves it automatically
using whatever package manager is available on your machine.

Modes:
  (default)    Interactive - shows plan, asks for confirmation
  --dry-run    Shows what would be done without executing
  --auto       Non-interactive mode for CI environments`,
	RunE: runFix,
}

var (
	fixDryRun bool
	fixAuto   bool
	fixJSON   bool
	fixConfig string
)

func init() {
	fixCmd.Flags().BoolVar(&fixDryRun, "dry-run", false, "Show fix plan without executing")
	fixCmd.Flags().BoolVar(&fixAuto, "auto", false, "Non-interactive mode (skip confirmation)")
	fixCmd.Flags().BoolVar(&fixJSON, "json", false, "Output fix plan as JSON")
	fixCmd.Flags().StringVarP(&fixConfig, "config", "c", "", "Path to .ground.yaml")
	rootCmd.AddCommand(fixCmd)
}

func runFix(cmd *cobra.Command, args []string) error {
	// Styles
	titleStyle := lipgloss.NewStyle().Bold(true)
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

	// Load config
	configPath := fixConfig
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

	// Detect tools and check drift
	names := make([]string, len(cfg.Tools))
	for i, t := range cfg.Tools {
		names[i] = t.Name
	}
	detected := detector.DetectByNames(names)
	report := drift.Compare(cfg, detected)

	if report.Summary.Errors == 0 && report.Summary.Warnings == 0 {
		fmt.Println()
		fmt.Printf("  %s All tools match the team standard. Nothing to fix.\n\n", okStyle.Render("[ok]"))
		return nil
	}

	// Detect package managers
	managers := pkgmanager.Detect()
	if len(managers) == 0 && !fixDryRun {
		fmt.Println()
		fmt.Printf("  %s No supported package manager found.\n", errStyle.Render("[ERR]"))
		fmt.Println("  Install one of: brew, apt, winget, scoop, choco, dnf, pacman")
		fmt.Println()
		return fmt.Errorf("no package manager available")
	}

	// Build fix plans
	plans := fixer.BuildFixPlans(report, managers)

	if len(plans) == 0 {
		fmt.Println()
		fmt.Printf("  %s Nothing to fix.\n\n", okStyle.Render("[ok]"))
		return nil
	}

	// JSON output
	if fixJSON {
		data, err := json.MarshalIndent(plans, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	// Display fix plan
	fmt.Println()
	fmt.Printf("%s %s\n", titleStyle.Render("groundctl"), dimStyle.Render("fix plan"))
	if len(managers) > 0 {
		pmNames := make([]string, len(managers))
		for i, m := range managers {
			pmNames[i] = m.Name
		}
		fmt.Printf("%s %s\n", dimStyle.Render("using:"), strings.Join(pmNames, ", "))
	}
	fmt.Println()

	hasManual := false
	for _, plan := range plans {
		actionLabel := "install"
		if plan.Action == fixer.ActionUpgrade {
			actionLabel = "upgrade"
		}

		if len(plan.Command) > 0 {
			icon := errStyle.Render("+")
			if plan.Action == fixer.ActionUpgrade {
				icon = warnStyle.Render("~")
			}
			fmt.Printf("  %s %-18s %s\n", icon, plan.Tool, cmdStyle.Render(plan.CommandStr))
			if plan.Current != "" {
				fmt.Printf("    %s %s -> %s\n", dimStyle.Render(actionLabel+":"), plan.Current, plan.Expected)
			}
		} else {
			fmt.Printf("  %s %-18s %s\n", dimStyle.Render("?"), plan.Tool, dimStyle.Render("(manual) "+plan.ManualHint))
			hasManual = true
		}
	}

	autoCount := 0
	for _, p := range plans {
		if len(p.Command) > 0 {
			autoCount++
		}
	}

	fmt.Println()
	fmt.Printf("  %s\n", dimStyle.Render(fmt.Sprintf("%d fixes planned (%d automatic, %d manual)",
		len(plans), autoCount, len(plans)-autoCount)))

	// Dry run stops here
	if fixDryRun {
		fmt.Println()
		fmt.Printf("  %s\n\n", dimStyle.Render("Dry run - no changes made."))
		return nil
	}

	// Confirm
	if !fixAuto {
		fmt.Println()
		fmt.Print("  Apply fixes? [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println()
			fmt.Printf("  %s\n\n", dimStyle.Render("Aborted."))
			return nil
		}
	}

	// Execute
	fmt.Println()
	fmt.Printf("  %s\n\n", titleStyle.Render("Applying fixes..."))

	execPlans := make([]fixer.FixPlan, 0, autoCount)
	for _, p := range plans {
		if len(p.Command) > 0 {
			execPlans = append(execPlans, p)
		}
	}

	results, allSuccess := fixer.ExecuteAll(execPlans, true)

	for _, r := range results {
		if r.Success {
			fmt.Printf("  %s %s\n", okStyle.Render("[ok]"), r.Plan.Tool)
		} else {
			fmt.Printf("  %s %s - %s\n", errStyle.Render("[FAIL]"), r.Plan.Tool, r.Error)
			if r.Output != "" {
				for _, line := range strings.Split(r.Output, "\n") {
					fmt.Printf("    %s\n", dimStyle.Render(line))
				}
			}
		}
	}

	fmt.Println()
	if allSuccess {
		fmt.Printf("  %s\n", okStyle.Render("All fixes applied successfully."))
	} else {
		fmt.Printf("  %s\n", warnStyle.Render("Some fixes failed. Review the output above."))
	}

	if hasManual {
		fmt.Println()
		fmt.Printf("  %s\n", dimStyle.Render("Some tools require manual installation - see hints above."))
	}

	fmt.Println()
	return nil
}
