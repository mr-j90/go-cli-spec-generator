package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-spec",
	Short: "AI-powered spec generator for Go projects",
	Long:  `go-spec is a CLI tool that uses AI to generate specs based on user inputs.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
