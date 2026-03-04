package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/groundctl/groundctl/internal/config"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate .ground.yaml syntax and constraints",
	Long:  `Checks .ground.yaml for common issues: missing fields, invalid versions, duplicate tools.`,
	RunE:  runValidate,
}

var (
	validateJSON   bool
	validateConfig string
)

func init() {
	validateCmd.Flags().BoolVar(&validateJSON, "json", false, "Output validation errors as JSON")
	validateCmd.Flags().StringVarP(&validateConfig, "config", "c", "", "Path to .ground.yaml")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	configPath := validateConfig
	if configPath == "" {
		configPath = findConfig()
	}
	if configPath == "" {
		configPath = filepath.Join(".", config.DefaultConfigFile)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	result := config.Validate(cfg)

	if validateJSON {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if result.IsValid() {
		fmt.Printf("%s %s is valid (%d tools)\n",
			okStyle.Render("✓"), configPath, len(cfg.Tools))
		return nil
	}

	fmt.Println()
	fmt.Printf("%s %s has %d issue(s):\n\n",
		errStyle.Render("✗"), configPath, len(result.Errors))

	for _, e := range result.Errors {
		fmt.Printf("  %s %s: %s\n",
			errStyle.Render("•"),
			dimStyle.Render(e.Field),
			e.Message)
	}
	fmt.Println()

	return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
}
