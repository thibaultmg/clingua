package oxford

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/language"

	"github.com/rs/zerolog/log"

	"github.com/thibaultmg/clingua/internal/entity"
	languageuc "github.com/thibaultmg/clingua/internal/usecase/language"
)

type Repo struct {
	baseURL *url.URL
	appID   string
	appKey  string
	client  *http.Client
	toLang  language.Tag
}

func New(client *http.Client, baseURL, appID, appKey string, to language.Tag) (Repo, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return Repo{}, err
	}

	return Repo{
		baseURL: u,
		appID:   appID,
		appKey:  appKey,
		client:  client,
		toLang:  to,
	}, nil
}

//nolint:funlen
func (r Repo) GetDefinition(ctx context.Context, word string, pos entity.PartOfSpeech) ([]languageuc.DefinitionEntry, error) {
	var ret []languageuc.DefinitionEntry

	var reqURL strings.Builder

	baseQueryString := "fields=definitions,domains,examples,pronunciations,registers&strictMatch=false"

	if _, err := fmt.Fprintf(&reqURL, "/api/v2/entries/%s/%s?%s", r.toLang, word, baseQueryString); err != nil {
		panic(err)
	}

	if !pos.IsAny() {
		if _, err := fmt.Fprintf(&reqURL, "&lexicalCategory=%s", pos); err != nil {
			panic(err)
		}
	}

	uriRef, err := url.Parse(reqURL.String())
	if err != nil {
		return ret, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL.ResolveReference(uriRef).String(), nil)
	if err != nil {
		return ret, err
	}

	req.Header.Add("app_id", r.appID)
	req.Header.Add("app_key", r.appKey)
	req.Header.Add("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return ret, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)

		_, err := buf.ReadFrom(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read get definition response body")
		}

		return ret, fmt.Errorf("http request %v failed with status %v and message: %s", req.URL, resp.StatusCode, buf.String())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}

	// uncomment for debug purpose
	// err = os.WriteFile("response.json", []byte(body), 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var entries EntriesResponse
	if err := json.Unmarshal(body, &entries); err != nil {
		return ret, err
	}

	ret = response2Internal(entries)

	return ret, nil
}
