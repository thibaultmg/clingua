package oxford

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/thibaultmg/clingua/internal/card"
	"github.com/thibaultmg/clingua/internal/entity"
	"golang.org/x/text/language"
)

type Repo struct {
	baseUrl *url.URL
	appID   string
	appKey  string
	client  *http.Client
}

func New(baseUrl, appID, appKey string) (Repo, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return Repo{}, err
	}

	return Repo{
		baseUrl: u,
		appID:   appID,
		appKey:  appKey,
		client:  &http.Client{},
	}, nil
}

func (r Repo) Get(ctx context.Context, word string, lang language.Tag, pos entity.PartOfSpeech) ([]card.DefinitionEntry, error) {
	var ret []card.DefinitionEntry

	queryString := fmt.Sprintf("fields=definitions,domains,examples,pronunciations,registers&strictMatch=false")
	if !pos.IsAny() {
		queryString = queryString + fmt.Sprintf("&lexicalCategory=%s", pos)
	}

	uriRef, err := url.Parse(fmt.Sprintf("/entries/%s/%s?%s", lang, word, queryString))
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

	body, err := io.ReadAll(resp.Body)
	var entries EntriesResponse
	if err := json.Unmarshal(body, &entries); err != nil {
		return ret, err
	}

	for _, r := range entries.Results {
		for _, le := range r.LexicalEntries {
			for _, e := range le.Entries {
				for _, s := range e.Senses {
					newEntry := card.DefinitionEntry{
						Definition:   s.Definitions[0],
						PartOfSpeech: entity.ParsePartOfSpeech(le.LexicalCategory.Text),
					}

					for _, ex := range s.Examples {
						newEntry.Exemples = append(newEntry.Exemples, ex.Text)
					}

					for _, reg := range s.Registers {
						newEntry.Registers = append(newEntry.Registers, reg.Text)
					}

					// for _, sub := range s.Subsenses {
					// 	newEntry.Domains = append(newEntry.Domains, sub.Domains)
					// }

					ret = append(ret, newEntry)

				}
			}
		}
	}

	return ret, nil
}
