package cmd

import (
	"fmt"

	"github.com/Ravenium22/groundctl/internal/shell"
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook [shell]",
	Short: "Generate shell hooks for auto-checking on cd",
	Long: `Outputs a shell script that auto-runs 'ground check' when you cd into
a directory containing .ground.yaml.

Usage:
  eval "$(ground hook bash)"    # Add to ~/.bashrc
  eval "$(ground hook zsh)"     # Add to ~/.zshrc
  ground hook fish | source     # Add to ~/.config/fish/config.fish
  ground hook --prompt starship # Prompt segment for Starship`,
	Args: func(cmd *cobra.Command, args []string) error {
		if hookPrompt != "" {
			if hookPrompt != "starship" && hookPrompt != "p10k" {
				return fmt.Errorf("invalid --prompt value %q (use: starship, p10k)", hookPrompt)
			}
			if len(args) != 0 {
				return fmt.Errorf("do not pass a shell argument when using --prompt")
			}
			return nil
		}

		if len(args) != 1 {
			return fmt.Errorf("requires exactly 1 shell argument (bash, zsh, fish, powershell)")
		}
		return nil
	},
	ValidArgs: []string{"bash", "zsh", "fish", "powershell", "pwsh"},
	RunE:      runHook,
}

var hookPrompt string

func init() {
	hookCmd.Flags().StringVar(&hookPrompt, "prompt", "", "Prompt integration snippet: starship or p10k")
	rootCmd.AddCommand(hookCmd)
}

func runHook(cmd *cobra.Command, args []string) error {
	if hookPrompt != "" {
		switch hookPrompt {
		case "starship":
			fmt.Print(shell.StarshipSegment())
		case "p10k":
			fmt.Print(shell.P10kSegment())
		default:
			return fmt.Errorf("prompt integration available for: starship, p10k")
		}
		return nil
	}

	shellName := args[0]

	switch shellName {
	case "bash":
		fmt.Print(shell.BashHook())
	case "zsh":
		fmt.Print(shell.ZshHook())
	case "fish":
		fmt.Print(shell.FishHook())
	case "powershell", "pwsh":
		fmt.Print(shell.PowerShellHook())
	default:
		return fmt.Errorf("unsupported shell %q. Supported: bash, zsh, fish, powershell", shellName)
	}
	return nil
}
