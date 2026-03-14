package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-spec",
	Short: "An AI-powered spec generator",
	Long:  `go-spec is a CLI tool that uses AI to generate specs based on user inputs.`,
}

// Execute runs the root command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
