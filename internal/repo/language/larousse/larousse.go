package larousse

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/rs/zerolog/log"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/usecase/language"
)

const baseURL = "https://www.larousse.fr/dictionnaires/anglais-francais/"

type HTTPDoClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Larousse repo returns single words translations.
type Repo struct {
	BaseURL    *url.URL
	httpClient HTTPDoClient
}

func New(httpClient HTTPDoClient) *Repo {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic("invalid base url")
	}

	return &Repo{
		httpClient: httpClient,
		BaseURL:    u,
	}
}

//nolint:lll
func (l *Repo) TranslateWord(ctx context.Context, rawWord string, pos entity.PartOfSpeech) ([]language.WordTranslationEntry, error) {
	word := strings.ReplaceAll(strings.TrimSpace(strings.ToLower(rawWord)), " ", "_")

	wordURL := *l.BaseURL
	wordURL.Path = path.Join(wordURL.Path, word)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wordURL.String(), nil)
	if err != nil {
		return []language.WordTranslationEntry{}, err
	}

	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"95\", \"Chromium\";v=\"95\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.larousse.fr/dictionnaires/anglais-francais/woolly/624783")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7,es;q=0.6")
	req.Header.Set("Cookie", "OptanonAlertBoxClosed=2021-11-09T10:15:54.155Z; eupubconsent-v2=CPPZudpPPZudpAcABBENB0CgAAAAAAAAAChQAAAAAAAA.YAAAAAAAAAAA; ASP.NET_SessionId=me4luinth2ivtet2gpbwqeq5; OptanonConsent=isGpcEnabled=0&datestamp=Tue+Nov+09+2021+11%3A46%3A02+GMT%2B0100+(heure+normale+d%E2%80%99Europe+centrale)&version=6.24.0&isIABGlobal=false&hosts=&consentId=aba48c65-d25d-46ae-a681-34d664e547a3&interactionCount=1&landingPath=NotLandingPage&groups=1%3A1%2C3%3A0%2C5%3A0%2C2%3A0%2C4%3A0&geolocation=FR%3BOCC&AwaitingReconsent=false")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return []language.WordTranslationEntry{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []language.WordTranslationEntry{}, errors.New("invalid http response")
	}

	respData, err := parseHTMLData(resp.Body)
	if err != nil {
		return []language.WordTranslationEntry{}, err
	}

	return response2WordTranslations(respData), nil
}

func parseHTMLData(body io.Reader) (responseData, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return responseData{}, err
	}

	articleNode := getArticleNode(doc)
	if articleNode == nil {
		return responseData{}, errors.New("failed to get article node")
	}

	ret := responseData{}

	var (
		parseData func(*html.Node)
		curPos    *partOfSpeech
		curItem   *posItem
		curEx     *example
	)

	addItem := func() {
		curEx = nil

		curPos.Items = append(curPos.Items, posItem{})
		curItem = &curPos.Items[len(curPos.Items)-1]
	}

	parseData = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// New pos item
			if attrSliceContains(n.Attr, atom.Class.String(), "ZoneEntree") {
				curItem = nil
				curEx = nil

				ret.Data = append(ret.Data, partOfSpeech{})
				curPos = &ret.Data[len(ret.Data)-1]
			}

			// New example
			if attrSliceContains(n.Attr, atom.Class.String(), "ZoneExpression") {
				if curItem != nil {
					curItem.Examples = append(curItem.Examples, example{})
					curEx = &curItem.Examples[len(curItem.Examples)-1]
				}
			}

			// Set pos
			if attrSliceContains(n.Attr, atom.Class.String(), "CategorieGrammaticale") {
				if curPos != nil {
					curPos.PartOfSpeech = getPos(n)
				}
			}

			// Set short definition
			if attrSliceContains(n.Attr, atom.Class.String(), "Indicateur") {
				if curItem != nil {
					curItem.Definition = getMeaning(n)
				} else if curPos != nil {
					addItem()
					curItem.Definition = getMeaning(n)
				}
			}

			// Add translation
			if attrSliceContains(n.Attr, atom.Class.String(), "Traduction") {
				if curItem != nil {
					curItem.Translations = append(curItem.Translations, getTranslation(n))
				} else if curPos != nil {
					addItem()
					curItem.Translations = append(curItem.Translations, getTranslation(n))
				}
			}

			if attrSliceContains(n.Attr, atom.Class.String(), "Locution2") {
				if curEx != nil {
					curEx.Example = renderNode(n)
				}
			}

			if attrSliceContains(n.Attr, atom.Class.String(), "Traduction2") {
				if curEx != nil {
					curEx.Translation = renderNode(n)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseData(c)
		}
	}
	parseData(articleNode)

	return ret, nil
}

func response2WordTranslations(res responseData) []language.WordTranslationEntry {
	ret := []language.WordTranslationEntry{}

	for _, pos := range res.Data {
		posVal, err := entity.ParsePartOfSpeech(pos.PartOfSpeech)
		if err != nil {
			if strings.Contains(pos.PartOfSpeech, "verb") {
				posVal = entity.Verb
			} else {
				log.Debug().Msgf("invalid pos value %s", pos.PartOfSpeech)

				continue
			}
		}

		for _, item := range pos.Items {
			for _, trans := range item.Translations {
				ret = append(ret, language.WordTranslationEntry{
					PartOfSpeech: posVal,
					Translation:  trans,
					Meaning:      item.Definition,
				})
			}
		}
	}

	return ret
}

func getArticleNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.DataAtom == atom.Div {
		if attrSliceContains(n.Attr, atom.Class.String(), "article_bilingue") {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res := getArticleNode(c)
		if res != nil {
			return res
		}
	}

	return nil
}

func getPos(n *html.Node) string {
	return strings.TrimSpace(stripHTMLTags(renderNode(n)))
}

var meaningRe = regexp.MustCompile(`\[(.+)\]`)

func getMeaning(n *html.Node) string {
	res := meaningRe.FindAllStringSubmatch(renderNode(n), -1)
	if len(res) > 0 && len(res[0]) > 1 {
		return res[0][1]
	}

	return ""
}

var fixSpacesRe = regexp.MustCompile(`\s+`)

func getTranslation(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if attrSliceContains(c.Attr, atom.Class.String(), "Genre") {
			n.RemoveChild(c)
		}

		if attrSliceContains(c.Attr, atom.Class.String(), "lienconj2") {
			n.RemoveChild(c)
		}
	}

	ret := stripHTMLTags(renderNode(n))

	return fixSpacesRe.ReplaceAllString(ret, " ")
}

func attrSliceContains(attr []html.Attribute, key, val string) bool {
	for _, e := range attr {
		if e.Key == key && strings.EqualFold(e.Val, val) {
			return true
		}
	}

	return false
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)

	if err := html.Render(w, n); err != nil {
		panic("failed to render node: " + err.Error())
	}

	return buf.String()
}
