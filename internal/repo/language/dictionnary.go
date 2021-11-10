package language

import (
	"log"
	"net/http"

	lang "golang.org/x/text/language"

	"github.com/thibaultmg/clingua/internal/config"
	"github.com/thibaultmg/clingua/internal/repo/language/deepl"
	"github.com/thibaultmg/clingua/internal/repo/language/larousse"
	"github.com/thibaultmg/clingua/internal/repo/language/oxford"
	languageuc "github.com/thibaultmg/clingua/internal/usecase/language"
)

// Factory for creating dictionary in the language.
func NewEnglishDictionnary(l lang.Tag) languageuc.Dictionary {
	httpClient := &http.Client{}

	var (
		ret languageuc.Dictionary
		err error
	)

	baseURL, appID, appKey := config.GetOxfordRepo()

	ret, err = oxford.New(httpClient, baseURL, appID, appKey, l)
	if err != nil {
		log.Fatalf("Error creating dictionary repo: %v", err)
	}

	return ret
}

// Factory for creating translator in the language.
func NewTranslator(fromLang lang.Tag, toLang lang.Tag) languageuc.Translator {
	httpClient := &http.Client{}

	var (
		ret languageuc.Translator
		err error
	)

	baseURL, authKey := config.GetDeeplRepo()

	ret, err = deepl.New(httpClient, authKey, baseURL, fromLang, toLang)
	if err != nil {
		log.Fatalf("Error creating translator repo: %v", err)
	}

	return ret
}

// Factory for creating word translator in the language.
func NewWordTranslator(fromLang lang.Tag, toLang lang.Tag) languageuc.WordTranslator {
	httpClient := &http.Client{}

	return larousse.New(httpClient)
}
