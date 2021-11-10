package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/presenter/cli/card"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vocabulary card",
	Long:  `Interactive command to create a vocabulary card, selecting definition and examples`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := entity.NewCard()
		c.Title = strings.Join(args, " ")

		cardEditor := makeCardEditor()
		cardEditor.SetCard(&c)

		cardPresenter := card.NewCardCLI(cardEditor)
		cardPresenter.RunCreate()
	},
}
