package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/Ravenium22/groundctl/internal/config"
	"github.com/Ravenium22/groundctl/internal/pkgmanager"
	"github.com/Ravenium22/groundctl/internal/profile"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose groundctl configuration and environment",
	Long:  `Runs a series of checks to verify that groundctl is properly configured and your environment is healthy.`,
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) error {
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	fmt.Println()
	fmt.Printf("  %s %s/%s, Go %s\n",
		dimStyle.Render("system:"), runtime.GOOS, runtime.GOARCH, runtime.Version())
	fmt.Printf("  %s %s (%s, %s)\n",
		dimStyle.Render("groundctl:"), Version, Commit, Date)
	fmt.Println()

	issues := 0

	// Check 1: .ground.yaml exists
	configPath := findConfig()
	if configPath != "" {
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("  %s .ground.yaml found but invalid: %v\n", errStyle.Render("[ERR]"), err)
			issues++
		} else {
			result := config.Validate(cfg)
			if result.IsValid() {
				fmt.Printf("  %s .ground.yaml valid (%d tools) at %s\n",
					okStyle.Render("[ok]"), len(cfg.Tools), configPath)
			} else {
				fmt.Printf("  %s .ground.yaml has %d issue(s)\n",
					warnStyle.Render("[!!]"), len(result.Errors))
				for _, e := range result.Errors {
					fmt.Printf("       %s\n", dimStyle.Render(e.Error()))
				}
				issues++
			}
		}
	} else {
		fmt.Printf("  %s No .ground.yaml found in current directory or parents\n",
			warnStyle.Render("[!!]"))
		issues++
	}

	// Check 2: Git available
	if _, err := exec.LookPath("git"); err == nil {
		fmt.Printf("  %s git is available\n", okStyle.Render("[ok]"))
	} else {
		fmt.Printf("  %s git not found (needed for ground pull/push)\n",
			errStyle.Render("[ERR]"))
		issues++
	}

	// Check 3: Package managers
	managers := pkgmanager.Detect()
	if len(managers) > 0 {
		names := ""
		for i, m := range managers {
			if i > 0 {
				names += ", "
			}
			names += m.Name
		}
		fmt.Printf("  %s Package manager(s): %s\n", okStyle.Render("[ok]"), names)
	} else {
		fmt.Printf("  %s No package managers found (ground fix won't work)\n",
			warnStyle.Render("[!!]"))
		issues++
	}

	// Check 4: Profiles directory
	profDir, err := profile.Dir()
	if err == nil {
		if _, err := os.Stat(profDir); err == nil {
			profiles, _ := profile.List()
			fmt.Printf("  %s Profiles directory exists (%d profiles)\n",
				okStyle.Render("[ok]"), len(profiles))
		} else {
			fmt.Printf("  %s Profiles directory not created yet %s\n",
				dimStyle.Render("[--]"),
				dimStyle.Render("(will be created on first 'ground profile save')"))
		}
	}

	// Check 5: Active profile
	active := profile.GetActive()
	if active != "" {
		fmt.Printf("  %s Active profile: %s\n", okStyle.Render("[ok]"), active)
	} else {
		fmt.Printf("  %s No active profile set\n", dimStyle.Render("[--]"))
	}

	// Check 6: Shell
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		shellEnv = os.Getenv("ComSpec")
	}
	if shellEnv != "" {
		fmt.Printf("  %s Shell: %s\n", okStyle.Render("[ok]"), filepath.Base(shellEnv))
	}

	// Summary
	fmt.Println()
	if issues == 0 {
		fmt.Printf("  %s\n", okStyle.Render("All checks passed. groundctl is ready."))
	} else {
		fmt.Printf("  %s\n", warnStyle.Render(fmt.Sprintf("%d issue(s) found. See above for details.", issues)))
	}
	fmt.Println()

	return nil
}
