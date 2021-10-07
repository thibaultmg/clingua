package oxford

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/usecase/card"
	"golang.org/x/text/language"
)

type Repo struct {
	baseUrl *url.URL
	appID   string
	appKey  string
	client  *http.Client
}

func New(client *http.Client, baseUrl, appID, appKey string) (Repo, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return Repo{}, err
	}

	return Repo{
		baseUrl: u,
		appID:   appID,
		appKey:  appKey,
		client:  client,
	}, nil
}

func (r Repo) GetDefinition(ctx context.Context, word string, lang language.Tag, pos entity.PartOfSpeech) ([]card.DefinitionEntry, error) {
	var ret []card.DefinitionEntry

	var reqUrl strings.Builder

	baseQueryString := "fields=definitions,domains,examples,pronunciations,registers&strictMatch=false"

	if _, err := fmt.Fprintf(&reqUrl, "/entries/%s/%s?%s", lang, word, baseQueryString); err != nil {
		panic(err)
	}

	if !pos.IsAny() {
		if _, err := fmt.Fprintf(&reqUrl, "&lexicalCategory=%s", pos); err != nil {
			panic(err)
		}
	}

	uriRef, err := url.Parse(reqUrl.String())
	if err != nil {
		return ret, err
	}

	req, err := http.NewRequest(http.MethodGet, r.baseUrl.ResolveReference(uriRef).String(), nil)
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
		return ret, fmt.Errorf("http request failed with status %v and message: %s", resp.StatusCode, resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}

	var entries EntriesResponse
	if err := json.Unmarshal(body, &entries); err != nil {
		return ret, err
	}

	ret = response2Internal(entries)

	return ret, nil
}
