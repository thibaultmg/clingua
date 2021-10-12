package language

import (
	"log"
	"net/http"

	"github.com/thibaultmg/clingua/internal/config"
	"github.com/thibaultmg/clingua/internal/repo/language/deepl"
	"github.com/thibaultmg/clingua/internal/repo/language/oxford"
	languageuc "github.com/thibaultmg/clingua/internal/usecase/language"
	lang "golang.org/x/text/language"
)

// Factory for creating dictionnary in the language
func NewDictionnary(toLanguage lang.Tag) languageuc.Dictionnary {
	httpClient := &http.Client{}

	var ret languageuc.Dictionnary
	var err error

	if toLanguage == lang.English {
		baseUrl, appID, appKey := config.GetOxfordRepo()
		ret, err = oxford.New(httpClient, baseUrl, appID, appKey, toLanguage)
		if err != nil {
			log.Fatalf("Error creating oxford repo: %v", err)
		}

	}

	return ret
}

// Factory for creating translator in the language
func NewTranslator(fromLang lang.Tag, toLang lang.Tag) languageuc.Translator {
	httpClient := &http.Client{}

	var ret languageuc.Translator
	var err error

	if toLang == lang.English || toLang == lang.French {
		baseUrl, authKey := config.GetDeeplRepo()
		ret, err = deepl.New(httpClient, authKey, baseUrl, fromLang, toLang)
		if err != nil {
			log.Fatalf("Error creating oxford repo: %v", err)
		}

	} else {
		panic("no implemented")
	}

	return ret
}
