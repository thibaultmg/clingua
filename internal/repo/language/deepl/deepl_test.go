package deepl_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/thibaultmg/clingua/internal/repo/language/deepl"
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
	t.Parallel()

	word := "ace"
	authKey := "myauthkey"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test url path
		assert.True(strings.HasPrefix(req.URL.String(), "/v2/translate"))

		// Test auth key
		authHeader := req.Header.Get("Authorization")
		assert.Contains(authHeader, authKey)

		// Test payload
		payload, err := io.ReadAll(req.Body)
		assert.Nil(err)
		assert.Contains(string(payload), "text="+word)

		// Send mocked response
		_, writeErr := rw.Write([]byte(testResponse))
		assert.Nil(writeErr)
	}))

	defer server.Close()

	// Test repo with mocked server
	deeplRep, err := deepl.New(server.Client(), authKey, server.URL, language.French, language.English)
	assert.Nil(err)

	res, err := deeplRep.Translate(context.Background(), word)
	assert.Nil(err)
	assert.Len(res, 1)
	assert.Equal("Hallo, Welt!", res[0])
}
