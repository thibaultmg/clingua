package deepl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/thibaultmg/clingua/internal/entity"
	"golang.org/x/text/language"
)

// curl https://api-free.deepl.com/v2/translate \
// 	-d auth_key=5106829b-9ffd-8118-1cbe-0c75a5bb9b2d:fx \
// 	-d "text=Hello, world!"  \
// 	-d "target_lang=DE"

// {
// 	"translations": [{
// 		"detected_source_language":"EN",
// 		"text":"Hallo, Welt!"
// 	}]
// }

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
	baseUrl  *url.URL
	fromLang language.Tag
	toLang   language.Tag
}

func New(client *http.Client, authKey string, baseUrl string, from, to language.Tag) (Repo, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return Repo{}, err
	}

	return Repo{
		client:   client,
		authKey:  authKey,
		baseUrl:  u,
		fromLang: from,
		toLang:   to,
	}, nil
}

func (r Repo) GetTranslation(ctx context.Context, text string, pos entity.PartOfSpeech) ([]string, error) {
	var ret []string

	reqUrl := fmt.Sprintf("/v2/translate?auth_key=%s&text=%s&source_lang=%s&target_lang=%s&split_sentences=0", r.authKey, text, r.fromLang, r.toLang)
	uriRef, err := url.Parse(reqUrl)
	if err != nil {
		return ret, err
	}

	req, err := http.NewRequest(http.MethodGet, r.baseUrl.ResolveReference(uriRef).String(), nil)
	if err != nil {
		return ret, err
	}

	fmt.Println(req.URL.String())

	req.Header.Add("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return ret, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return ret, fmt.Errorf("http request %v failed with status %v and message: %s", req.URL, resp.StatusCode, buf.String())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}

	fmt.Println(string(body))

	var entries Response
	if err := json.Unmarshal(body, &entries); err != nil {
		return ret, err
	}

	ret = make([]string, 0, len(entries.Translations))
	for _, tradText := range entries.Translations {
		ret = append(ret, tradText.Text)
	}

	return ret, nil
}
