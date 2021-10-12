package deepl_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/repo/language/deepl"
	"golang.org/x/text/language"
)

var testResponse = `
{
	"translations": [{
		"detected_source_language":"EN",
		"text":"Hallo, Welt!"
	}]
}
`

func TestDeepl_Nominal(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	word := "ace"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test url path
		assert.True(strings.HasPrefix(req.URL.String(), "/v2/translate"))

		// Send mocked response
		_, writeErr := rw.Write([]byte(testResponse))
		assert.Nil(writeErr)
	}))

	defer server.Close()

	// Test repo with mocked server
	deeplRep, err := deepl.New(server.Client(), "authkey", server.URL, language.French, language.English)
	assert.Nil(err)

	res, err := deeplRep.GetTranslation(context.Background(), word, entity.Any)
	assert.Nil(err)
	assert.Len(res, 1)
	assert.Equal("Hallo, Welt!", res[0])
}
