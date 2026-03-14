package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zyx-holdings/go-spec/internal/tui"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Start an interactive spec generation session",
	Long:  `Launch the interactive TUI to generate a new specification document.`,
	RunE:  runNew,
}

func runNew(cmd *cobra.Command, args []string) error {
	if !isTerminal() {
		fmt.Fprintln(cmd.ErrOrStderr(), "error: specgen requires an interactive terminal")
		return fmt.Errorf("not an interactive terminal")
	}

	_, err := tui.Run(noColor)
	return err
}

// isTerminal reports whether os.Stdin is an interactive terminal.
func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func init() {
	rootCmd.AddCommand(newCmd)
}
