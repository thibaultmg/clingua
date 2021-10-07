package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var ClinguaVersion = "development"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Clingua",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ClinguaVersion)
	},
}
