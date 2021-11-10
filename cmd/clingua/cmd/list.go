package cmd

import (
	"github.com/spf13/cobra"

	"github.com/thibaultmg/clingua/internal/presenter/cli/card"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List and search vocabulary cards",
	Long:  `Interactive command to list vocabulary cards and editing them`,
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cardEditor := makeCardEditor()
		cardPresenter := card.NewCardCLI(cardEditor)

		cardPresenter.RunList()
	},
}
