package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/Ravenium22/groundctl/internal/config"
	"github.com/Ravenium22/groundctl/internal/team"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull <source>",
	Short: "Fetch a team standard from a git repo, URL, or local path",
	Long: `Fetches a .ground.yaml from a remote source and saves it locally.

Sources:
  github.com/org/repo     Clone git repo and extract .ground.yaml
  https://.../.ground.yaml  Download raw YAML file
  ./path/to/.ground.yaml  Copy from local filesystem`,
	Args: cobra.ExactArgs(1),
	RunE: runPull,
}

var pullOutput string

func init() {
	pullCmd.Flags().StringVarP(&pullOutput, "output", "o", "", "Output path (default: .ground.yaml)")
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	source := args[0]
	dest := pullOutput
	if dest == "" {
		dest = filepath.Join(".", config.DefaultConfigFile)
	}

	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	fmt.Printf("%s %s\n", dimStyle.Render("Pulling from"), source)

	cfg, err := team.Pull(source, dest)
	if err != nil {
		return fmt.Errorf("pull failed: %w", err)
	}

	fmt.Printf("%s Saved %s (%d tools)\n",
		okStyle.Render("[ok]"), dest, len(cfg.Tools))

	if cfg.Name != "" {
		fmt.Printf("  %s %s\n", dimStyle.Render("profile:"), cfg.Name)
	}

	return nil
}
