package strategy

import (
	"net/http"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type call func(int) (*http.Response, error)

type IStrategy interface {
	Apply(call) string
}

type baseStrategy struct {
}

func (baseStrategy) sanitizeHtmlToTweet(htmlToSanitize string) string {
	p := bluemonday.StrictPolicy()
	htmlSanitized := p.Sanitize(htmlToSanitize)

	fields := strings.Fields(string(htmlSanitized))
	htmlTextCleaned := strings.Join(fields, " ")

	const textMaxLen = 280
	if len(htmlTextCleaned) > textMaxLen {
		htmlTextCleaned = htmlTextCleaned[:textMaxLen]
	}

	return htmlTextCleaned
}
