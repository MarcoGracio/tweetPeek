package strategy

import (
	"errors"
	"net/http"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type requestProcessor func(int) (*http.Response, error)

type IStrategy interface {
	Apply(requestProcessor) string
	GetCurrentAttempt() int
}

type baseStrategy struct {
	maxRetries     int
	currentAttempt int
}

func NewBaseStrategy(maxRetries int) (*baseStrategy, error) {
	if maxRetries <= 0 {
		return nil, errors.New("max retries must be greater than 0")
	}
	return &baseStrategy{maxRetries: maxRetries}, nil
}

func (*baseStrategy) sanitizeHtmlToTweet(htmlToSanitize string) string {
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

func (strategy *baseStrategy) GetCurrentAttempt() int {
	return strategy.currentAttempt
}
