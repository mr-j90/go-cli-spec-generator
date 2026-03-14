package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Exit codes as defined in the spec.
const (
	ExitSuccess    = 0
	ExitValidation = 1
	ExitExport     = 2

	Version = "0.1.0"
)

var noColor bool

var rootCmd = &cobra.Command{
	Use:     "specgen",
	Short:   "AI-powered spec generator CLI",
	Long:    `specgen is an interactive CLI tool for generating structured specification documents.`,
	Version: Version,
}

// Execute runs the root command and exits with the appropriate code on failure.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitValidation)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
}
