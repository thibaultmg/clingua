package deepl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/text/language"
)

type Response struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	SourceLang string `json:"detected_source_language"`
	Text       string `json:"text"`
}

type Repo struct {
	client   *http.Client
	authKey  string
	baseURL  *url.URL
	fromLang language.Tag
	toLang   language.Tag
}

func New(client *http.Client, authKey string, baseURL string, from, to language.Tag) (Repo, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return Repo{}, err
	}

	return Repo{
		client:   client,
		authKey:  authKey,
		baseURL:  u,
		fromLang: from,
		toLang:   to,
	}, nil
}

func (r Repo) Translate(ctx context.Context, text string) ([]string, error) {
	reqURL := fmt.Sprintf("/v2/translate?auth_key=%s&text=%s&source_lang=%s&target_lang=%s&split_sentences=0",
		r.authKey, text, r.fromLang, r.toLang)

	uriRef, err := url.Parse(reqURL)
	if err != nil {
		return []string{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL.ResolveReference(uriRef).String(), nil)
	if err != nil {
		return []string{}, err
	}

	req.Header.Add("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return []string{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)

		if _, err := buf.ReadFrom(resp.Body); err != nil {
			return []string{}, err
		}

		return []string{}, fmt.Errorf("http request %v failed with status %v and message: %s", req.URL, resp.StatusCode, buf.String())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var entries Response
	if err := json.Unmarshal(body, &entries); err != nil {
		return []string{}, err
	}

	ret := make([]string, 0, len(entries.Translations))
	for _, tradText := range entries.Translations {
		ret = append(ret, tradText.Text)
	}

	return ret, nil
}
