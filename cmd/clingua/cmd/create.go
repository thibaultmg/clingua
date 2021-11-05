package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thibaultmg/clingua/internal/config"
	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/presenter/cli/card"
	"github.com/thibaultmg/clingua/internal/repo/language"
	"github.com/thibaultmg/clingua/internal/repo/store/filesystem"
	carduc "github.com/thibaultmg/clingua/internal/usecase/card"
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
		fromLang, toLang := config.GetLanguages()
		dict := language.NewDictionnary(toLang)
		trans := language.NewTranslator(toLang, fromLang)
		luc := languageuc.New(dict, trans)

		cardRepo := filesystem.New(config.GetFSRepoPath())
		cuc := carduc.New(cardRepo)

		c := entity.NewCard()
		c.Title = strings.Join(args, " ")
		cardEditor := card.NewCardEditor(&c, luc, cuc)

		cardPresenter := card.NewCardCLI(cardEditor)

		cardPresenter.Run()
		fmt.Println("bye bye")
	},
}
