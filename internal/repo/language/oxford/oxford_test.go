package oxford_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/repo/language/oxford"
)

var testResponse = `{"id":"ace","metadata":{"operation":"retrieve","provider":"Oxford University Press","schema":"RetrieveEntry"},"results":[{"id":"ace","language":"en-gb","lexicalEntries":[{"entries":[{"homographNumber":"102","pronunciations":[{"audioFile":"https://audio.oxforddictionaries.com/en/mp3/ace_1_gb_1_abbr.mp3","dialects":["British English"],"phoneticNotation":"IPA","phoneticSpelling":"eɪs"}],"senses":[{"definitions":["(in tennis and similar games) serve an ace against (an opponent)"],"examples":[{"text":"he can ace opponents with serves of no more than 62 mph"}],"id":"m_en_gbus0005680.020","registers":[{"id":"informal","text":"Informal"}],"subsenses":[{"definitions":["score an ace on (a hole) or with (a shot)"],"domains":[{"id":"golf","text":"Golf"}],"examples":[{"text":"there was a prize for the first player to ace the hole"}],"id":"m_en_gbus0005680.026"}]},{"definitions":["achieve high marks in (a test or exam)"],"examples":[{"text":"I aced my grammar test"}],"id":"m_en_gbus0005680.028","registers":[{"id":"informal","text":"Informal"}],"subsenses":[{"definitions":["outdo someone in a competitive situation"],"examples":[{"text":"the magazine won an award, acing out its rivals"}],"id":"m_en_gbus0005680.029"}]}]}],"language":"en-gb","lexicalCategory":{"id":"verb","text":"Verb"},"text":"ace"}],"type":"headword","word":"ace"}],"word":"ace"}`

func TestOxford_Nominal(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	word := "ace"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test url path
		assert.True(strings.HasPrefix(req.URL.String(), "/entries/en/"+word))

		// Send mocked response
		_, writeErr := rw.Write([]byte(testResponse))
		assert.Nil(writeErr)
	}))

	defer server.Close()

	// Test repo with mocked server
	oxRep, err := oxford.New(server.Client(), server.URL, "appid", "appkey")
	assert.Nil(err)

	res, err := oxRep.Get(context.Background(), word, language.English, entity.Any)
	assert.Nil(err)
	assert.Len(res, 2)
}
