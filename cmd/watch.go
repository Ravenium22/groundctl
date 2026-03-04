package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/groundctl/groundctl/internal/config"
	"github.com/groundctl/groundctl/internal/detector"
	"github.com/groundctl/groundctl/internal/drift"
	"github.com/groundctl/groundctl/internal/model"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for drift changes in the background",
	Long: `Runs periodic drift checks and reports when your environment drifts
from the team standard. Press Ctrl+C to stop.`,
	RunE: runWatch,
}

var (
	watchInterval time.Duration
	watchConfig   string
)

func init() {
	watchCmd.Flags().DurationVarP(&watchInterval, "interval", "i", 30*time.Second, "Check interval (e.g. 30s, 1m, 5m)")
	watchCmd.Flags().StringVarP(&watchConfig, "config", "c", "", "Path to .ground.yaml")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	configPath := watchConfig
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

	names := make([]string, len(cfg.Tools))
	for i, t := range cfg.Tools {
		names[i] = t.Name
	}

	fmt.Printf("%s Watching for drift every %s (Ctrl+C to stop)\n",
		dimStyle.Render("[watch]"), watchInterval)
	fmt.Printf("%s config: %s\n\n", dimStyle.Render("[watch]"), configPath)

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(watchInterval)
	defer ticker.Stop()

	// Run initial check
	checkAndReport(cfg, names, dimStyle, okStyle, warnStyle, errStyle)

	for {
		select {
		case <-ticker.C:
			// Reload config in case it changed
			freshCfg, err := config.Load(configPath)
			if err != nil {
				fmt.Printf("%s Could not reload config: %v\n",
					errStyle.Render("[watch]"), err)
				continue
			}
			freshNames := make([]string, len(freshCfg.Tools))
			for i, t := range freshCfg.Tools {
				freshNames[i] = t.Name
			}
			checkAndReport(freshCfg, freshNames, dimStyle, okStyle, warnStyle, errStyle)
		case <-sigCh:
			fmt.Printf("\n%s Stopped watching.\n", dimStyle.Render("[watch]"))
			return nil
		}
	}
}

func checkAndReport(cfg *model.GroundConfig, names []string, dimStyle, okStyle, warnStyle, errStyle lipgloss.Style) {
	detected := detector.DetectByNames(names)
	report := drift.Compare(cfg, detected)

	ts := time.Now().Format("15:04:05")

	if report.Summary.Errors == 0 && report.Summary.Warnings == 0 {
		fmt.Printf("%s %s — all %d tools ok\n",
			dimStyle.Render("["+ts+"]"),
			okStyle.Render("clean"),
			report.Summary.Total)
	} else {
		parts := []string{}
		if report.Summary.Errors > 0 {
			parts = append(parts, errStyle.Render(fmt.Sprintf("%d errors", report.Summary.Errors)))
		}
		if report.Summary.Warnings > 0 {
			parts = append(parts, warnStyle.Render(fmt.Sprintf("%d warnings", report.Summary.Warnings)))
		}
		msg := ""
		for i, p := range parts {
			if i > 0 {
				msg += ", "
			}
			msg += p
		}
		fmt.Printf("%s %s — %s\n",
			dimStyle.Render("["+ts+"]"),
			errStyle.Render("drift"),
			msg)
	}
}
