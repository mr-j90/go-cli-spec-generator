package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var resumeSession string

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume an interrupted spec generation session",
	Long:  `Resume a previously interrupted interactive spec generation session from a saved session file.`,
	RunE:  runResume,
}

func runResume(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(cmd.OutOrStdout(), "Resuming spec generation session...")
	if resumeSession != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Session: %s\n", resumeSession)
	}
	// TODO: restore session via internal/session and launch TUI
	return nil
}

func init() {
	rootCmd.AddCommand(resumeCmd)
	resumeCmd.Flags().StringVar(&resumeSession, "session", "", "session file path to resume")
}
