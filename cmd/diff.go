package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/Ravenium22/groundctl/internal/model"
	"github.com/Ravenium22/groundctl/internal/snapshot"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show what changed since last snapshot",
	Long:  `Compares the current machine state against the last saved snapshot and shows differences.`,
	RunE:  runDiff,
}

var (
	diffSnapshot string
	diffJSON     bool
)

func init() {
	diffCmd.Flags().StringVarP(&diffSnapshot, "from", "f", "", "Path to previous snapshot JSON")
	diffCmd.Flags().BoolVar(&diffJSON, "json", false, "Output diff as JSON")
	rootCmd.AddCommand(diffCmd)
}

// DiffItem represents a single change between snapshots.
type DiffItem struct {
	Tool    string `json:"tool"`
	Change  string `json:"change"` // "added", "removed", "upgraded", "downgraded"
	Before  string `json:"before,omitempty"`
	After   string `json:"after,omitempty"`
}

func runDiff(cmd *cobra.Command, args []string) error {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	removeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	changeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

	// Load previous snapshot
	snapPath := diffSnapshot
	if snapPath == "" {
		home, _ := os.UserHomeDir()
		snapPath = filepath.Join(home, ".groundctl", "last-snapshot.json")
	}

	prevData, err := os.ReadFile(snapPath)
	if err != nil {
		return fmt.Errorf("no previous snapshot found at %s. Run 'ground snapshot -o %s' first", snapPath, snapPath)
	}

	var prev model.Snapshot
	if err := json.Unmarshal(prevData, &prev); err != nil {
		return fmt.Errorf("could not parse previous snapshot: %w", err)
	}

	// Capture current state
	current, err := snapshot.Capture()
	if err != nil {
		return fmt.Errorf("could not capture current state: %w", err)
	}

	// Build diff
	prevMap := make(map[string]model.DetectedTool)
	for _, t := range prev.Tools {
		prevMap[t.Name] = t
	}
	currMap := make(map[string]model.DetectedTool)
	for _, t := range current.Tools {
		currMap[t.Name] = t
	}

	var diffs []DiffItem

	// Check for changes and additions
	for _, curr := range current.Tools {
		old, existed := prevMap[curr.Name]
		if !existed || (!old.Found && curr.Found) {
			if curr.Found {
				diffs = append(diffs, DiffItem{
					Tool: curr.Name, Change: "added", After: curr.Version,
				})
			}
			continue
		}
		if old.Found && !curr.Found {
			diffs = append(diffs, DiffItem{
				Tool: curr.Name, Change: "removed", Before: old.Version,
			})
			continue
		}
		if old.Version != curr.Version && old.Found && curr.Found {
			change := "changed"
			diffs = append(diffs, DiffItem{
				Tool: curr.Name, Change: change, Before: old.Version, After: curr.Version,
			})
		}
	}

	// Check for removals
	for _, old := range prev.Tools {
		if _, exists := currMap[old.Name]; !exists && old.Found {
			diffs = append(diffs, DiffItem{
				Tool: old.Name, Change: "removed", Before: old.Version,
			})
		}
	}

	if diffJSON {
		data, _ := json.MarshalIndent(diffs, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if len(diffs) == 0 {
		fmt.Printf("  %s No changes since last snapshot.\n", dimStyle.Render("[ok]"))
		return nil
	}

	fmt.Println()
	fmt.Printf("  %s\n\n", dimStyle.Render(fmt.Sprintf("%d change(s) since last snapshot:", len(diffs))))

	for _, d := range diffs {
		switch d.Change {
		case "added":
			fmt.Printf("  %s %-18s %s\n", addStyle.Render("+"), d.Tool, addStyle.Render(d.After))
		case "removed":
			fmt.Printf("  %s %-18s %s\n", removeStyle.Render("-"), d.Tool, removeStyle.Render(d.Before))
		default:
			fmt.Printf("  %s %-18s %s -> %s\n", changeStyle.Render("~"), d.Tool,
				dimStyle.Render(d.Before), changeStyle.Render(d.After))
		}
	}
	fmt.Println()

	// Save current as new baseline
	newData, _ := snapshot.ToJSON(current)
	_ = os.MkdirAll(filepath.Dir(snapPath), 0755)
	_ = os.WriteFile(snapPath, newData, 0644)
	fmt.Printf("  %s\n", dimStyle.Render("Snapshot updated."))

	return nil
}
