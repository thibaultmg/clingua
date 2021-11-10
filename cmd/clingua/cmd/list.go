package cmd

import (
	"github.com/spf13/cobra"

	"github.com/thibaultmg/clingua/internal/config"
	"github.com/thibaultmg/clingua/internal/presenter/cli/card"
	"github.com/thibaultmg/clingua/internal/repo/language"
	"github.com/thibaultmg/clingua/internal/repo/store/filesystem"
	carduc "github.com/thibaultmg/clingua/internal/usecase/card"
	languageuc "github.com/thibaultmg/clingua/internal/usecase/language"
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
		fromLang, toLang := config.GetLanguages()
		dict := language.NewEnglishDictionnary(toLang)
		trans := language.NewTranslator(toLang, fromLang)
		wtrans := language.NewWordTranslator(toLang, fromLang)
		luc := languageuc.New(dict, trans, wtrans)

		cardRepo := filesystem.New(config.GetFSRepoPath())
		cuc := carduc.New(cardRepo)

		cardEditor := card.NewCardEditor(nil, luc, cuc)
		cardPresenter := card.NewCardCLI(cardEditor)

		cardPresenter.RunList()
	},
}
