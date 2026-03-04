package report

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Ravenium22/groundctl/internal/model"
)

// FormatJSON returns the drift report as indented JSON.
func FormatJSON(report *model.DriftReport) (string, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal report: %w", err)
	}
	return string(data), nil
}

// FormatMarkdown returns the drift report as a GitHub-flavored Markdown table.
func FormatMarkdown(report *model.DriftReport) string {
	var b strings.Builder

	b.WriteString("# groundctl Drift Report\n\n")
	b.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().UTC().Format(time.RFC3339)))

	b.WriteString("## Summary\n\n")
	b.WriteString(fmt.Sprintf("| Metric | Count |\n"))
	b.WriteString(fmt.Sprintf("| --- | --- |\n"))
	b.WriteString(fmt.Sprintf("| Total | %d |\n", report.Summary.Total))
	b.WriteString(fmt.Sprintf("| OK | %d |\n", report.Summary.OK))
	b.WriteString(fmt.Sprintf("| Warnings | %d |\n", report.Summary.Warnings))
	b.WriteString(fmt.Sprintf("| Errors | %d |\n", report.Summary.Errors))

	b.WriteString("\n## Tools\n\n")
	b.WriteString("| Tool | Status | Expected | Actual | Message |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")

	for _, item := range report.Items {
		statusIcon := "?"
		switch item.Status {
		case model.DriftOK:
			statusIcon = "ok"
		case model.DriftWarning:
			statusIcon = "warning"
		case model.DriftError:
			statusIcon = "error"
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			item.Tool, statusIcon,
			orDash(item.Expected), orDash(item.Actual), orDash(item.Message)))
	}

	return b.String()
}

// FormatHTML returns the drift report as a standalone HTML page.
func FormatHTML(report *model.DriftReport) string {
	var b strings.Builder

	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>groundctl Drift Report</title>
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 800px; margin: 40px auto; padding: 0 20px; color: #1a1a2e; background: #f8f9fa; }
  h1 { color: #16213e; border-bottom: 2px solid #0f3460; padding-bottom: 8px; }
  h2 { color: #16213e; margin-top: 32px; }
  .summary { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin: 16px 0; }
  .stat { background: white; border-radius: 8px; padding: 16px; text-align: center; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
  .stat .num { font-size: 28px; font-weight: bold; }
  .stat .label { font-size: 13px; color: #666; margin-top: 4px; }
  .ok .num { color: #22c55e; }
  .warn .num { color: #eab308; }
  .err .num { color: #ef4444; }
  table { width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
  th { background: #16213e; color: white; padding: 10px 12px; text-align: left; font-size: 13px; }
  td { padding: 8px 12px; border-bottom: 1px solid #eee; font-size: 14px; }
  tr:last-child td { border-bottom: none; }
  .status-ok { color: #22c55e; font-weight: bold; }
  .status-warning { color: #eab308; font-weight: bold; }
  .status-error { color: #ef4444; font-weight: bold; }
  .ts { font-size: 13px; color: #888; margin-top: 24px; }
</style>
</head>
<body>
<h1>groundctl Drift Report</h1>
`)

	b.WriteString(`<div class="summary">`)
	b.WriteString(fmt.Sprintf(`<div class="stat"><div class="num">%d</div><div class="label">Total</div></div>`, report.Summary.Total))
	b.WriteString(fmt.Sprintf(`<div class="stat ok"><div class="num">%d</div><div class="label">OK</div></div>`, report.Summary.OK))
	b.WriteString(fmt.Sprintf(`<div class="stat warn"><div class="num">%d</div><div class="label">Warnings</div></div>`, report.Summary.Warnings))
	b.WriteString(fmt.Sprintf(`<div class="stat err"><div class="num">%d</div><div class="label">Errors</div></div>`, report.Summary.Errors))
	b.WriteString(`</div>`)

	b.WriteString(`<h2>Tools</h2>
<table>
<tr><th>Tool</th><th>Status</th><th>Expected</th><th>Actual</th><th>Message</th></tr>
`)

	for _, item := range report.Items {
		statusClass := "status-ok"
		statusText := "ok"
		switch item.Status {
		case model.DriftWarning:
			statusClass = "status-warning"
			statusText = "warning"
		case model.DriftError:
			statusClass = "status-error"
			statusText = "error"
		}
		b.WriteString(fmt.Sprintf(`<tr><td><strong>%s</strong></td><td class="%s">%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`,
			item.Tool, statusClass, statusText,
			orDash(item.Expected), orDash(item.Actual), orDash(item.Message)))
		b.WriteString("\n")
	}

	b.WriteString(`</table>`)
	b.WriteString(fmt.Sprintf(`<p class="ts">Generated: %s</p>`, time.Now().UTC().Format(time.RFC3339)))
	b.WriteString("\n</body>\n</html>\n")

	return b.String()
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
