package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number and exit",
	Long:  `Print the version number of specgen and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "specgen version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
