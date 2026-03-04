package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/groundctl/groundctl/internal/recovery"
	"github.com/groundctl/groundctl/internal/telemetry"
	"github.com/spf13/cobra"
)

var (
	Verbose     bool
	Debug       bool
	NoTelemetry bool
)

var rootCmd = &cobra.Command{
	Use:   "ground",
	Short: "groundctl - terraform plan for your local dev machine",
	Long: `groundctl detects how your local development environment has drifted
from your team's standard and fixes it with one command.

  ground init       Create a .ground.yaml from your current machine
  ground snapshot   Capture current machine state as JSON
  ground check      Compare your machine against the team standard
  ground fix        Auto-fix detected drift
  ground profile    Manage environment profiles
  ground doctor     Diagnose groundctl configuration`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if NoTelemetry {
			os.Setenv("GROUND_TELEMETRY", "off")
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug mode (show detection commands and timing)")
	rootCmd.PersistentFlags().BoolVar(&NoTelemetry, "no-telemetry", false, "Disable telemetry for this invocation")
}

func Execute() {
	defer recovery.Wrap("root")()
	start := time.Now()
	err := rootCmd.Execute()
	duration := time.Since(start)

	// Record telemetry
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	telemetry.Record(telemetry.Event{
		Command:    rootCmd.CalledAs(),
		DurationMs: duration.Milliseconds(),
		ExitCode:   exitCode,
		Version:    Version,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
