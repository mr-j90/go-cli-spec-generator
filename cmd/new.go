package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Start an interactive spec generation session",
	Long:  `Launch the interactive TUI to generate a new specification document.`,
	RunE:  runNew,
}

func runNew(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(cmd.OutOrStdout(), "Launching interactive spec generator...")
	// TODO: launch TUI via internal/tui
	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)
}
