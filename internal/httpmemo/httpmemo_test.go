package httpmemo_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thibaultmg/clingua/internal/httpmemo"
)

func TestHTTPMemo_Get(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	payload := []byte("coucou")
	counter := 0

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, err := rw.Write(payload)
		assert.Nil(err)
		counter++
	}))

	memo := httpmemo.New(server.Client().Do)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil) //nolint:noctx
	assert.Nil(err)

	// Do a first request, check that server is hit and response returned
	res, err := memo.Do(req)
	assert.Nil(err)

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.Nil(err)
	assert.Equal(payload, data)
	assert.Equal(1, counter)

	// Do a second request, check that server is not hit and response returned
	resCache, err := memo.Do(req)
	assert.Nil(err)

	defer resCache.Body.Close()

	data, err = io.ReadAll(resCache.Body)
	assert.Nil(err)
	assert.Equal(payload, data)
	assert.Equal(1, counter)
}

func TestHTTPMemo_Post(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	payload := "coucou"
	counter := 0

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		counter++
	}))

	memo := httpmemo.New(server.Client().Do)

	req, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(payload)) //nolint:noctx
	assert.Nil(err)

	// Do a first request, check that server is hit and response returned
	res, err := memo.Do(req)
	assert.Nil(err)

	defer res.Body.Close()

	assert.Equal(1, counter)

	// Do a second request, change payload, check that server is hit again
	req, err = http.NewRequest(http.MethodPost, server.URL, strings.NewReader("new payload")) //nolint:noctx
	assert.Nil(err)

	resCache, err := memo.Do(req)
	assert.Nil(err)

	defer resCache.Body.Close()

	assert.Equal(2, counter)
}
