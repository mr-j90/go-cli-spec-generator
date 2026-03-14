package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
