package larousse

import (
	"html"
	"regexp"
	"strings"
)

const regex = `<.*?>`

func stripHTMLTags(s string) string {
	r := regexp.MustCompile(regex)
	s = r.ReplaceAllString(s, "")

	return strings.TrimSpace(html.UnescapeString(s))
}
