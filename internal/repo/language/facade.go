package language

import (
	"context"
	"log"
	"net/http"

	lang "golang.org/x/text/language"

	"github.com/thibaultmg/clingua/internal/config"
	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/repo/language/oxford"
	"github.com/thibaultmg/clingua/internal/usecase/card"
)

// Language struct acts as a facade integrating all language external services into a simplified interface.
// A language instance is dedicated to one fromLanguage (the learner native language) and one toLanguage
// (the learned language).
type Language struct {
	fromLanguage lang.Tag
	toLanguage   lang.Tag
	dictionnary  Dictionnary
}

func New(fromLanguage, toLanguage lang.Tag) *Language {
	ret := &Language{}

	httpClient := &http.Client{}

	if toLanguage == lang.English {
		baseUrl, appID, appKey := config.GetOxfordRepo()
		dict, err := oxford.New(httpClient, baseUrl, appID, appKey)
		if err != nil {
			log.Fatalf("Error creating oxford repo: %v", err)
		}

		ret.dictionnary = dict
	}
	return ret
}

func (l *Language) Define(ctx context.Context, word string, pos entity.PartOfSpeech) ([]card.DefinitionEntry, error) {
	return l.dictionnary.GetDefinition(ctx, word, l.toLanguage, pos)
}
