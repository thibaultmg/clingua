package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/thibaultmg/clingua/internal/config"
	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/presenter/cli/card"
	"github.com/thibaultmg/clingua/internal/repo/language"
	languageuc "github.com/thibaultmg/clingua/internal/usecase/language"
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
		_, toLanguage := config.GetLanguages()
		dict := language.NewDictionnary(toLanguage)
		luc := languageuc.New(dict)
		cardPresenter := card.New(luc)

		cardPresenter.CreateCard(strings.Join(args, " "), entity.Any)
	},
}
