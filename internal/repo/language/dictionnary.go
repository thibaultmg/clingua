package language

import (
	"log"
	"net/http"

	lang "golang.org/x/text/language"

	"github.com/thibaultmg/clingua/internal/config"
	"github.com/thibaultmg/clingua/internal/httpmemo"
	"github.com/thibaultmg/clingua/internal/repo/language/deepl"
	"github.com/thibaultmg/clingua/internal/repo/language/larousse"
	"github.com/thibaultmg/clingua/internal/repo/language/oxford"
)

// Factory for creating dictionary in the language.
func NewEnglishDictionnary(l lang.Tag) *oxford.Repo {
	httpClient := &http.Client{}
	memoClient := httpmemo.New(httpClient.Do)

	baseURL, appID, appKey := config.GetOxfordRepo()

	ret, err := oxford.New(memoClient, baseURL, appID, appKey, l)
	if err != nil {
		log.Fatalf("Error creating dictionary repo: %v", err)
	}

	return ret
}

// Factory for creating translator in the language.
func NewTranslator(fromLang lang.Tag, toLang lang.Tag) *deepl.Repo {
	httpClient := &http.Client{}
	memoClient := httpmemo.New(httpClient.Do)

	baseURL, authKey := config.GetDeeplRepo()

	ret, err := deepl.New(memoClient, authKey, baseURL, fromLang, toLang)
	if err != nil {
		log.Fatalf("Error creating translator repo: %v", err)
	}

	return ret
}

// Factory for creating word translator in the language.
func NewWordTranslator(fromLang lang.Tag, toLang lang.Tag) *larousse.Repo {
	httpClient := &http.Client{}
	memoClient := httpmemo.New(httpClient.Do)

	return larousse.New(memoClient)
}
