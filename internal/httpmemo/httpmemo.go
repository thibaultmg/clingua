package httpmemo

import (
	"bytes"
	"crypto/sha1" // nolint:gosec
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

const initCacheSize = 10

type HTTPDoFunc func(req *http.Request) (*http.Response, error)

type httpRes struct {
	res  *http.Response
	body []byte
	err  error
}

type HTTPMemo struct {
	cache map[string]httpRes
	do    HTTPDoFunc
}

func New(do HTTPDoFunc) *HTTPMemo {
	return &HTTPMemo{
		cache: make(map[string]httpRes, initCacheSize),
		do:    do,
	}
}

func (h *HTTPMemo) Do(req *http.Request) (*http.Response, error) {
	reqHash := makeRequestHash(req)

	cacheVal, ok := h.cache[reqHash]
	if !ok {
		var netErr net.Error

		cacheVal.res, cacheVal.err = h.do(req)
		if cacheVal.err != nil && errors.As(cacheVal.err, &netErr) {
			// Do not cahce if error is temporary
			if netErr.Temporary() {
				return cacheVal.res, cacheVal.err
			}
		}

		body, err := io.ReadAll(cacheVal.res.Body)
		if err != nil {
			log.Warn().Err(err).Msg("failed to read http response body")

			return cacheVal.res, cacheVal.err
		}

		cacheVal.body = body

		h.cache[reqHash] = cacheVal
	}

	cacheVal.res.Body = io.NopCloser(bytes.NewReader(cacheVal.body))

	return cacheVal.res, cacheVal.err
}

func makeRequestHash(req *http.Request) string {
	var (
		body []byte
		err  error
	)

	if req.Body != nil {
		defer req.Body.Close()

		body, err = io.ReadAll(req.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read request body")
		}

		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	h := sha1.New() // nolint:gosec
	_, err = io.WriteString(h, req.Method)
	checkErr(err)

	_, err = io.WriteString(h, req.URL.String())
	checkErr(err)

	_, err = io.WriteString(h, string(body))
	checkErr(err)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func checkErr(err error) {
	if err != nil {
		log.Error().Err(err).Msg("failed to write string into hash writer")
	}
}
