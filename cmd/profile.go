package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/Ravenium22/groundctl/internal/config"
	"github.com/Ravenium22/groundctl/internal/profile"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage environment profiles",
	Long: `Manage multiple environment profiles (work, personal, client-X).

Profiles are stored in ~/.groundctl/profiles/ and can inherit from each other.`,
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	RunE:  runProfileList,
}

var profileSwitchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch to a named profile",
	Long:  `Applies a profile by resolving its inheritance chain and writing .ground.yaml.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileSwitch,
}

var profileSaveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save the current .ground.yaml as a named profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileSave,
}

var profileShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show a profile's resolved configuration",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileShow,
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a named profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileDelete,
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileSwitchCmd)
	profileCmd.AddCommand(profileSaveCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileDeleteCmd)
	rootCmd.AddCommand(profileCmd)
}

func runProfileList(cmd *cobra.Command, args []string) error {
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	names, err := profile.List()
	if err != nil {
		return err
	}

	active := profile.GetActive()

	if len(names) == 0 {
		fmt.Println(dimStyle.Render("No profiles saved. Use 'ground profile save <name>' to create one."))
		return nil
	}

	fmt.Println()
	for _, name := range names {
		marker := "  "
		if name == active {
			marker = okStyle.Render("* ")
		}
		fmt.Printf("%s%s\n", marker, name)
	}
	fmt.Println()
	return nil
}

func runProfileSwitch(cmd *cobra.Command, args []string) error {
	name := args[0]
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	if !profile.Exists(name) {
		return fmt.Errorf("profile %q not found. Run 'ground profile list' to see available profiles", name)
	}

	cfg, err := profile.Resolve(name)
	if err != nil {
		return err
	}

	dest := filepath.Join(".", config.DefaultConfigFile)
	if err := config.Save(dest, cfg); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}

	if err := profile.SetActive(name); err != nil {
		return fmt.Errorf("could not set active profile: %w", err)
	}

	fmt.Printf("%s Switched to profile %q (%d tools)\n",
		okStyle.Render("[ok]"), name, len(cfg.Tools))

	if cfg.Extends != "" {
		fmt.Printf("  %s extends %q\n", dimStyle.Render("inherits:"), cfg.Extends)
	}
	return nil
}

func runProfileSave(cmd *cobra.Command, args []string) error {
	name := args[0]
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)

	src := filepath.Join(".", config.DefaultConfigFile)
	cfg, err := config.Load(src)
	if err != nil {
		return fmt.Errorf("could not load %s: %w. Run 'ground init' first", config.DefaultConfigFile, err)
	}

	if err := profile.Save(name, cfg); err != nil {
		return err
	}

	fmt.Printf("%s Profile %q saved (%d tools)\n",
		okStyle.Render("[ok]"), name, len(cfg.Tools))
	return nil
}

func runProfileShow(cmd *cobra.Command, args []string) error {
	name := args[0]
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	cfg, err := profile.Resolve(name)
	if err != nil {
		return err
	}

	fmt.Println()
	if cfg.Name != "" {
		fmt.Printf("  name: %s\n", cfg.Name)
	}
	if cfg.Extends != "" {
		fmt.Printf("  extends: %s\n", cfg.Extends)
	}
	fmt.Printf("  tools: %d\n\n", len(cfg.Tools))

	for _, t := range cfg.Tools {
		ver := t.Version
		if ver == "" {
			ver = "*"
		}
		sev := string(t.Severity)
		if sev == "" {
			sev = "required"
		}
		fmt.Printf("  %-18s %s %s\n", t.Name, dimStyle.Render(ver), dimStyle.Render("["+sev+"]"))
	}
	fmt.Println()
	return nil
}

func runProfileDelete(cmd *cobra.Command, args []string) error {
	name := args[0]
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)

	if !profile.Exists(name) {
		return fmt.Errorf("profile %q not found", name)
	}

	if err := profile.Delete(name); err != nil {
		return err
	}

	// If this was the active profile, clear it
	if profile.GetActive() == name {
		_ = profile.ClearActive()
	}

	fmt.Printf("%s Profile %q deleted\n", okStyle.Render("[ok]"), name)
	return nil
}
