package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/Ravenium22/groundctl/internal/config"
	"github.com/Ravenium22/groundctl/internal/secrets"
	"github.com/spf13/cobra"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage secret references in your config",
	Long:  `Manage secret references defined in .ground.yaml. Supports 1Password (op), HashiCorp Vault, OS keychain, and environment variables.`,
}

var secretsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate that all secret references are resolvable",
	Long:  `Checks each secret reference in .ground.yaml to verify the backend is available and the secret exists.`,
	RunE:  runSecretsCheck,
}

var secretsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secret references in config",
	RunE:  runSecretsList,
}

var secretsEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Generate a .env file from secret references",
	Long:  `Resolves all secret references and writes them as NAME=value pairs to a .env file.`,
	RunE:  runSecretsEnv,
}

var (
	secretsConfig string
	secretsOutput string
)

func init() {
	secretsCmd.PersistentFlags().StringVarP(&secretsConfig, "config", "c", "", "Path to .ground.yaml")

	secretsEnvCmd.Flags().StringVarP(&secretsOutput, "output", "o", ".env", "Output path for .env file")

	secretsCmd.AddCommand(secretsCheckCmd)
	secretsCmd.AddCommand(secretsListCmd)
	secretsCmd.AddCommand(secretsEnvCmd)
	rootCmd.AddCommand(secretsCmd)
}

func loadSecretsConfig() ([]secrets.SecretRef, *secrets.Registry, error) {
	configPath := secretsConfig
	if configPath == "" {
		configPath = findConfig()
	}
	if configPath == "" {
		configPath = filepath.Join(".", config.DefaultConfigFile)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, nil, err
	}

	if len(cfg.Secrets) == 0 {
		return nil, nil, fmt.Errorf("no secrets defined in %s", configPath)
	}

	registry := secrets.DefaultRegistry()
	var refs []secrets.SecretRef

	for _, s := range cfg.Secrets {
		ref := secrets.ParseRef(s.Ref)
		if ref == nil {
			return nil, nil, fmt.Errorf("invalid secret reference %q for %s", s.Ref, s.Name)
		}
		refs = append(refs, *ref)
	}

	return refs, registry, nil
}

func runSecretsCheck(cmd *cobra.Command, args []string) error {
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	configPath := secretsConfig
	if configPath == "" {
		configPath = findConfig()
	}
	if configPath == "" {
		configPath = filepath.Join(".", config.DefaultConfigFile)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	if len(cfg.Secrets) == 0 {
		fmt.Printf("  %s No secrets defined in config.\\n", dimStyle.Render("[--]"))
		return nil
	}

	registry := secrets.DefaultRegistry()
	issues := 0

	fmt.Println()
	for _, s := range cfg.Secrets {
		ref := secrets.ParseRef(s.Ref)
		if ref == nil {
			fmt.Printf("  %s %-20s invalid reference: %s\n", errStyle.Render("[ERR]"), s.Name, s.Ref)
			issues++
			continue
		}

		if err := registry.CheckRef(*ref); err != nil {
			fmt.Printf("  %s %-20s %s\n", errStyle.Render("[ERR]"), s.Name, dimStyle.Render(err.Error()))
			issues++
		} else {
			fmt.Printf("  %s %-20s %s %s\n", okStyle.Render("[ok]"), s.Name,
				dimStyle.Render(ref.Backend+"://"), dimStyle.Render(ref.Path))
		}
	}

	fmt.Println()
	if issues == 0 {
		fmt.Printf("  %s\n", okStyle.Render(fmt.Sprintf("All %d secret(s) resolvable.", len(cfg.Secrets))))
	} else {
		fmt.Printf("  %s\n", errStyle.Render(fmt.Sprintf("%d of %d secret(s) failed.", issues, len(cfg.Secrets))))
	}
	fmt.Println()

	if issues > 0 {
		return fmt.Errorf("%d secret(s) could not be resolved", issues)
	}
	return nil
}

func runSecretsList(cmd *cobra.Command, args []string) error {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	configPath := secretsConfig
	if configPath == "" {
		configPath = findConfig()
	}
	if configPath == "" {
		configPath = filepath.Join(".", config.DefaultConfigFile)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	if len(cfg.Secrets) == 0 {
		fmt.Printf("  %s No secrets defined in config.\n", dimStyle.Render("[--]"))
		return nil
	}

	fmt.Println()
	for _, s := range cfg.Secrets {
		ref := secrets.ParseRef(s.Ref)
		backend := "?"
		if ref != nil {
			backend = ref.Backend
		}
		desc := ""
		if s.Description != "" {
			desc = dimStyle.Render(" - " + s.Description)
		}
		fmt.Printf("  %-20s %-8s %s%s\n", s.Name, dimStyle.Render("["+backend+"]"), s.Ref, desc)
	}
	fmt.Println()

	return nil
}

func runSecretsEnv(cmd *cobra.Command, args []string) error {
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	configPath := secretsConfig
	if configPath == "" {
		configPath = findConfig()
	}
	if configPath == "" {
		configPath = filepath.Join(".", config.DefaultConfigFile)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	if len(cfg.Secrets) == 0 {
		fmt.Printf("  %s No secrets defined in config.\n", dimStyle.Render("[--]"))
		return nil
	}

	registry := secrets.DefaultRegistry()
	var lines []string
	issues := 0

	for _, s := range cfg.Secrets {
		ref := secrets.ParseRef(s.Ref)
		if ref == nil {
			fmt.Printf("  %s %-20s invalid reference: %s\n", errStyle.Render("[ERR]"), s.Name, s.Ref)
			issues++
			continue
		}

		val, err := registry.Resolve(*ref)
		if err != nil {
			fmt.Printf("  %s %-20s %s\n", errStyle.Render("[ERR]"), s.Name, err.Error())
			issues++
			continue
		}

		if s.Description != "" {
			lines = append(lines, fmt.Sprintf("# %s", s.Description))
		}
		lines = append(lines, fmt.Sprintf("%s=%s", s.Name, quoteEnvValue(val)))
		fmt.Printf("  %s %-20s %s\n", okStyle.Render("[ok]"), s.Name, secrets.MaskValue(val))
	}

	if issues > 0 {
		return fmt.Errorf("%d secret(s) could not be resolved", issues)
	}

	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(secretsOutput, []byte(content), 0600); err != nil {
		return fmt.Errorf("could not write %s: %w", secretsOutput, err)
	}

	fmt.Printf("\n  %s\n\n", okStyle.Render(fmt.Sprintf("Written %d secret(s) to %s", len(cfg.Secrets), secretsOutput)))
	return nil
}

// quoteEnvValue wraps a value in quotes if it contains spaces or special chars.
func quoteEnvValue(val string) string {
	if strings.ContainsAny(val, " \t\n\"'#$\\") {
		escaped := strings.ReplaceAll(val, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		return "\"" + escaped + "\""
	}
	return val
}
