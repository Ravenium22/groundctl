package cmd

import (
	"fmt"
	"os"

	"github.com/groundctl/groundctl/internal/snapshot"
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Capture current machine state as JSON",
	Long:  `Runs all tool detectors and outputs a JSON snapshot of your machine's current state.`,
	RunE:  runSnapshot,
}

var snapshotOutput string

func init() {
	snapshotCmd.Flags().StringVarP(&snapshotOutput, "output", "o", "", "Write snapshot to file instead of stdout")
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	snap, err := snapshot.Capture()
	if err != nil {
		return fmt.Errorf("snapshot failed: %w", err)
	}

	data, err := snapshot.ToJSON(snap)
	if err != nil {
		return err
	}

	if snapshotOutput != "" {
		if err := os.WriteFile(snapshotOutput, data, 0644); err != nil {
			return fmt.Errorf("failed to write snapshot: %w", err)
		}
		fmt.Printf("Snapshot written to %s\n", snapshotOutput)
		return nil
	}

	fmt.Println(string(data))
	return nil
}
