package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/groundctl/groundctl/internal/config"
	"github.com/groundctl/groundctl/internal/detector"
	"github.com/groundctl/groundctl/internal/drift"
	"github.com/groundctl/groundctl/internal/report"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a drift report in various formats",
	Long: `Generates a drift report comparing your machine against .ground.yaml.

Supported formats:
  json      JSON (default)
  markdown  GitHub-flavored Markdown table
  html      Standalone styled HTML page`,
	RunE: runReport,
}

var (
	reportFormat string
	reportOutput string
	reportConfig string
)

func init() {
	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "json", "Output format: json, markdown, html")
	reportCmd.Flags().StringVarP(&reportOutput, "output", "o", "", "Write to file (default: stdout)")
	reportCmd.Flags().StringVarP(&reportConfig, "config", "c", "", "Path to .ground.yaml")
	rootCmd.AddCommand(reportCmd)
}

func runReport(cmd *cobra.Command, args []string) error {
	configPath := reportConfig
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
	detected := detector.DetectByNames(names)
	driftReport := drift.Compare(cfg, detected)

	var output string
	switch reportFormat {
	case "json":
		output, err = report.FormatJSON(driftReport)
		if err != nil {
			return err
		}
	case "markdown", "md":
		output = report.FormatMarkdown(driftReport)
	case "html":
		output = report.FormatHTML(driftReport)
	default:
		return fmt.Errorf("unsupported format %q (use: json, markdown, html)", reportFormat)
	}

	if reportOutput != "" {
		if err := os.MkdirAll(filepath.Dir(reportOutput), 0755); err != nil {
			return fmt.Errorf("could not create output directory: %w", err)
		}
		if err := os.WriteFile(reportOutput, []byte(output), 0644); err != nil {
			return fmt.Errorf("could not write report: %w", err)
		}
		fmt.Printf("Report written to %s\n", reportOutput)
	} else {
		fmt.Print(output)
	}

	return nil
}
