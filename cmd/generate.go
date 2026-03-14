package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	outputFile string
	specType   string
)

var generateCmd = &cobra.Command{
	Use:   "generate [description]",
	Short: "Generate a spec from a description",
	Long:  `Generate an AI-powered spec document from a natural language description.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path (default: stdout)")
	generateCmd.Flags().StringVarP(&specType, "type", "t", "feature", "spec type (feature, api, architecture)")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	description := args[0]
	fmt.Fprintf(cmd.OutOrStdout(), "Generating %s spec for: %s\n", specType, description)
	if outputFile != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Output will be written to: %s\n", outputFile)
	}
	return nil
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a spec",
	Long:  `Generate an AI-powered spec based on user inputs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called — feature coming soon")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
