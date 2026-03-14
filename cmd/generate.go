package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	generateInput  string
	generateOutput string
	generateFormat string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a spec document from an input file",
	Long:  `Generate a specification document from an existing input file without launching the interactive TUI.`,
	RunE:  runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(cmd.OutOrStdout(), "Generating spec document...")
	if generateInput != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Input: %s\n", generateInput)
	}
	if generateOutput != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Output: %s\n", generateOutput)
	}
	if generateFormat != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Format: %s\n", generateFormat)
	}
	// TODO: implement spec generation using internal/render and internal/export
	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&generateInput, "input", "i", "", "input session file path")
	generateCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "output file path")
	generateCmd.Flags().StringVarP(&generateFormat, "format", "f", "markdown", "output format (markdown, pdf, docx)")
}
