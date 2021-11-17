package larousse_test

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/repo/language/larousse"
)

//go:embed testdata/response_envy.html
var envyResponse []byte

const apiBasePath = "/dictionnaires/anglais-francais"

func TestLarousse_TranslateWord(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	word := "car"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request validity
		assert.Contains(req.URL.String(), word)
		assert.Equal(http.MethodGet, req.Method)

		// Send mocked response
		rw.WriteHeader(http.StatusOK)
		_, writeErr := rw.Write(envyResponse)
		assert.Nil(writeErr)
	}))

	defer server.Close()

	larousseRepo := larousse.New(server.Client())
	u, _ := url.Parse(server.URL + apiBasePath)
	larousseRepo.BaseURL = u

	res, err := larousseRepo.TranslateWord(context.Background(), word, entity.Adjective)
	assert.Nil(err)
	assert.Greater(len(res), 0)
	assert.Equal(res[0].PartOfSpeech, entity.Noun)
	assert.Contains(res[0].Meaning, "jealousy")
	assert.Contains(res[0].Translation, "envie")
}
