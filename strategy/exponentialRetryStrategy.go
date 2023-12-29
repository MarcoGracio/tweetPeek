package strategy

import (
	"errors"
	"io"
	"net/http"
	"time"
)

type exponentialRetryStrategy struct {
	*baseStrategy
	initialBackoffSeconds int
}

func NewExponentialRetryStrategy(maxRetries int, initialBackoffSeconds int) (*exponentialRetryStrategy, error) {
	if initialBackoffSeconds <= 0 {
		return nil, errors.New("not negative backoff allowed")
	}

	strategy, err := NewBaseStrategy(maxRetries)

	return &exponentialRetryStrategy{baseStrategy: strategy, initialBackoffSeconds: initialBackoffSeconds}, err
}

func (strategy *exponentialRetryStrategy) Apply(request requestProcessor) string {
	var resp *http.Response
	var err error

	waitTime := time.Duration(strategy.initialBackoffSeconds) * time.Second

	for strategy.currentAttempt = 1; strategy.currentAttempt <= strategy.maxRetries; strategy.currentAttempt++ {
		resp, err = request(strategy.currentAttempt)

		if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			responseBody, _ := io.ReadAll(resp.Body)
			tweet := strategy.sanitizeHtmlToTweet(string(responseBody))
			return "exponentialRetryStrategy" + "\n" + tweet
		}

		time.Sleep(waitTime)
		waitTime *= 2
	}

	if err != nil {
		return err.Error()
	} else {
		return "Failed to get a valid response."
	}
}
