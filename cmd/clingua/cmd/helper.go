package cmd

import (
	"github.com/thibaultmg/clingua/internal/config"
	"github.com/thibaultmg/clingua/internal/presenter/cli/card"
	"github.com/thibaultmg/clingua/internal/repo/language"
	"github.com/thibaultmg/clingua/internal/repo/store/filesystem"
	carduc "github.com/thibaultmg/clingua/internal/usecase/card"
	languageuc "github.com/thibaultmg/clingua/internal/usecase/language"
)

func makeCardEditor() *card.CardEditor {
	fromLang, toLang := config.GetLanguages()
	dict := language.NewEnglishDictionnary(toLang)
	trans := language.NewTranslator(toLang, fromLang)
	wtrans := language.NewWordTranslator(toLang, fromLang)
	luc := languageuc.New(dict, trans, wtrans)

	cardRepo := filesystem.New(config.GetFSRepoPath())
	cuc := carduc.New(cardRepo)

	return card.NewCardEditor(nil, luc, cuc)
}
