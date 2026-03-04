package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/groundctl/groundctl/internal/config"
	"github.com/groundctl/groundctl/internal/team"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Publish .ground.yaml to the git remote",
	Long:  `Stages, commits, and pushes your .ground.yaml to the current git remote.`,
	RunE:  runPush,
}

var pushMessage string

func init() {
	pushCmd.Flags().StringVarP(&pushMessage, "message", "m", "", "Commit message")
	rootCmd.AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) error {
	configPath := filepath.Join(".", config.DefaultConfigFile)
	if !config.Exists(configPath) {
		return fmt.Errorf("no %s found. Run 'ground init' first", config.DefaultConfigFile)
	}

	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	fmt.Printf("%s\n", dimStyle.Render("Pushing .ground.yaml to remote..."))

	if err := team.Push(configPath, pushMessage); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	fmt.Printf("%s %s pushed to remote\n", okStyle.Render("✓"), config.DefaultConfigFile)
	return nil
}
