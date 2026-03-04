package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var feedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Submit feedback about groundctl",
	Long: `Submit feedback, bug reports, or feature requests directly from the CLI.
Feedback is saved locally and can optionally be submitted to the groundctl team.`,
	RunE: runFeedback,
}

var (
	feedbackCategory string
	feedbackMessage  string
	feedbackEmail    string
)

func init() {
	feedbackCmd.Flags().StringVarP(&feedbackCategory, "category", "t", "", "Category: bug, feature, other")
	feedbackCmd.Flags().StringVarP(&feedbackMessage, "message", "m", "", "Feedback message")
	feedbackCmd.Flags().StringVar(&feedbackEmail, "email", "", "Optional contact email")
	rootCmd.AddCommand(feedbackCmd)
}

// FeedbackEntry represents a saved feedback item.
type FeedbackEntry struct {
	Category  string `json:"category"`
	Message   string `json:"message"`
	Email     string `json:"email,omitempty"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

func runFeedback(cmd *cobra.Command, args []string) error {
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	reader := bufio.NewReader(os.Stdin)

	// Collect category
	category := feedbackCategory
	if category == "" {
		fmt.Printf("  %s (bug / feature / other): ", dimStyle.Render("Category"))
		line, _ := reader.ReadString('\n')
		category = strings.TrimSpace(line)
	}
	if category != "bug" && category != "feature" && category != "other" {
		category = "other"
	}

	// Collect message
	message := feedbackMessage
	if message == "" {
		fmt.Printf("  %s ", dimStyle.Render("Feedback:"))
		line, _ := reader.ReadString('\n')
		message = strings.TrimSpace(line)
	}
	if message == "" {
		return fmt.Errorf("feedback message cannot be empty")
	}

	// Collect optional email
	email := feedbackEmail
	if email == "" {
		fmt.Printf("  %s ", dimStyle.Render("Email (optional, press Enter to skip):"))
		line, _ := reader.ReadString('\n')
		email = strings.TrimSpace(line)
	}

	entry := FeedbackEntry{
		Category:  category,
		Message:   message,
		Email:     email,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Version:   Version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Save locally
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}
	feedbackDir := filepath.Join(home, ".groundctl", "feedback")
	if err := os.MkdirAll(feedbackDir, 0755); err != nil {
		return fmt.Errorf("could not create feedback directory: %w", err)
	}

	filename := fmt.Sprintf("%s-%d.json", category, time.Now().UnixNano())
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	feedbackPath := filepath.Join(feedbackDir, filename)
	if err := os.WriteFile(feedbackPath, data, 0644); err != nil {
		return fmt.Errorf("could not save feedback: %w", err)
	}

	fmt.Println()
	fmt.Printf("  %s Feedback saved.\n", okStyle.Render("[ok]"))
	fmt.Printf("  %s\n", dimStyle.Render(feedbackPath))
	fmt.Printf("  %s\n", dimStyle.Render("Thank you! To submit to the team, visit github.com/Ravenium22/groundctl/issues"))
	fmt.Println()

	return nil
}
