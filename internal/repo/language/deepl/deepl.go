package deepl

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
)

type HTTPDoClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Response struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	SourceLang string `json:"detected_source_language"`
	Text       string `json:"text"`
}

type Repo struct {
	client   HTTPDoClient
	authKey  string
	baseURL  *url.URL
	fromLang language.Tag
	toLang   language.Tag
}

func New(client HTTPDoClient, authKey string, baseURL string, from, to language.Tag) (*Repo, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return &Repo{}, err
	}

	u.Path = "/v2/translate"

	return &Repo{
		client:   client,
		authKey:  authKey,
		baseURL:  u,
		fromLang: from,
		toLang:   to,
	}, nil
}

func (r *Repo) Translate(ctx context.Context, text string) ([]string, error) {
	params := url.Values{}
	params.Add("text", text)
	params.Add("target_lang", r.toLang.String())
	params.Add("source_lang", r.fromLang.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return []string{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "DeepL-Auth-Key "+r.authKey)

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
