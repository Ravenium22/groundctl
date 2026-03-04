package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/Ravenium22/groundctl/internal/telemetry"
	"github.com/spf13/cobra"
)

var telemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: "Manage anonymous usage telemetry",
	Long: `groundctl collects anonymous usage telemetry to improve the tool.
Telemetry is opt-in and collects only: command name, OS, architecture,
tool count, and execution time. No personal information is collected.

  ground telemetry on     Enable telemetry
  ground telemetry off    Disable telemetry
  ground telemetry status Show current state`,
}

var telemetryOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Enable anonymous telemetry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := telemetry.SetEnabled(true); err != nil {
			return err
		}
		okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		fmt.Printf("  %s Telemetry enabled. Thank you for helping improve groundctl!\n", okStyle.Render("[ok]"))
		return nil
	},
}

var telemetryOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Disable anonymous telemetry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := telemetry.SetEnabled(false); err != nil {
			return err
		}
		if err := telemetry.ClearEvents(); err != nil {
			return err
		}
		dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		fmt.Printf("  %s Telemetry disabled. All stored events cleared.\n", dimStyle.Render("[--]"))
		return nil
	},
}

var telemetryStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show telemetry status",
	RunE: func(cmd *cobra.Command, args []string) error {
		okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

		if telemetry.IsEnabled() {
			count := telemetry.EventCount()
			fmt.Printf("  %s Telemetry is enabled (%d events stored)\n", okStyle.Render("[on]"), count)
		} else {
			fmt.Printf("  %s Telemetry is disabled\n", dimStyle.Render("[off]"))
		}

		fmt.Println()
		fmt.Printf("  %s\n", dimStyle.Render("We collect: command name, OS, arch, tool count, duration."))
		fmt.Printf("  %s\n", dimStyle.Render("We never collect: file paths, tool versions, secrets, or PII."))
		return nil
	},
}

func init() {
	telemetryCmd.AddCommand(telemetryOnCmd)
	telemetryCmd.AddCommand(telemetryOffCmd)
	telemetryCmd.AddCommand(telemetryStatusCmd)
	rootCmd.AddCommand(telemetryCmd)
}
