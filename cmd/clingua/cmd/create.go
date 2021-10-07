package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thibaultmg/clingua/internal/presenter/cli/card"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vocabulary card",
	Long:  `Interactive command to create a vocabulary card, selecting definition and exemples`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("creating card for", strings.Join(args, " "))
		card.Start()
	},
}
