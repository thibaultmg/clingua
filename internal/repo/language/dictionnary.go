package language

import (
	"log"
	"net/http"

	lang "golang.org/x/text/language"

	"github.com/thibaultmg/clingua/internal/config"
	"github.com/thibaultmg/clingua/internal/repo/language/deepl"
	"github.com/thibaultmg/clingua/internal/repo/language/oxford"
	languageuc "github.com/thibaultmg/clingua/internal/usecase/language"
)

// Factory for creating dictionary in the language.
func NewDictionnary(toLanguage lang.Tag) languageuc.Dictionary {
	httpClient := &http.Client{}

	var (
		ret languageuc.Dictionary
		err error
	)

	if toLanguage == lang.English {
		baseURL, appID, appKey := config.GetOxfordRepo()

		ret, err = oxford.New(httpClient, baseURL, appID, appKey, toLanguage)
		if err != nil {
			log.Fatalf("Error creating oxford repo: %v", err)
		}
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

	if toLang == lang.English || toLang == lang.French {
		baseURL, authKey := config.GetDeeplRepo()

		ret, err = deepl.New(httpClient, authKey, baseURL, fromLang, toLang)
		if err != nil {
			log.Fatalf("Error creating oxford repo: %v", err)
		}
	} else {
		panic("no implemented")
	}

	return ret
}
