package oxford_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/repo/language/oxford"
)

//go:embed testdata/response.json
var testResponse []byte

func TestOxford_Nominal(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	word := "ace"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test url path
		assert.Contains(req.URL.String(), "/entries/en/"+word)

		// Send mocked response
		assert.True(json.Valid(testResponse))
		// var payload bytes.Buffer

		// err := json.Compact(&payload, testResponse)
		// assert.Nil(err)

		_, writeErr := rw.Write(testResponse)
		assert.Nil(writeErr)
	}))

	defer server.Close()

	// Test repo with mocked server
	oxRep, err := oxford.New(server.Client(), server.URL, "appid", "appkey", language.English)
	assert.Nil(err)

	res, err := oxRep.GetDefinition(context.Background(), word, entity.Any)
	assert.Nil(err)
	assert.Len(res, 2)
}
