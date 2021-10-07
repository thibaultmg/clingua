package language

import (
	"log"
	"net/http"

	"github.com/thibaultmg/clingua/internal/config"
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
