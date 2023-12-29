package strategy

import (
	"io"
	"net/http"
)

type simpleRetryStrategy struct {
	*baseStrategy
}

func NewSimpleRetryStrategy(maxRetries int) (*simpleRetryStrategy, error) {
	strategy, err := NewBaseStrategy(maxRetries)
	return &simpleRetryStrategy{baseStrategy: strategy}, err
}

func (strategy *simpleRetryStrategy) Apply(processRequest requestProcessor) string {
	var resp *http.Response
	var err error
	for strategy.currentAttempt = 1; strategy.currentAttempt <= strategy.maxRetries; strategy.currentAttempt++ {
		resp, err = processRequest(strategy.currentAttempt)

		if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			responseBody, _ := io.ReadAll(resp.Body)
			tweet := strategy.sanitizeHtmlToTweet(string(responseBody))
			return "simpleRetryStrategy" + "\n" + tweet
		}
	}

	if err != nil {
		return err.Error()
	} else {
		return "Failed to get a valid response."
	}
}
