package resilience

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"tweetPeek/textProcessor"

	"github.com/microcosm-cc/bluemonday"
)

type requestProcessor func(int) (*http.Response, error)

type IStrategy interface {
	Apply(requestProcessor) (textProcessor.Tweets, error)
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

func (*baseStrategy) sanitizeBodyToTweets(body io.ReadCloser) textProcessor.Tweets {
	responseBody, _ := io.ReadAll(body)

	p := bluemonday.StrictPolicy()
	htmlSanitized := p.Sanitize(string(responseBody))

	fields := strings.Fields(string(htmlSanitized))
	htmlTextCleaned := strings.Join(fields, " ")

	return textProcessor.Tweets{}.Get(htmlTextCleaned, 3)
}

func (strategy *baseStrategy) GetCurrentAttempt() int {
	return strategy.currentAttempt
}
