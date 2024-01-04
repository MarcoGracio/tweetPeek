package resilience

import (
	"errors"
	"net/http"
	"time"
	"tweetPeek/textProcessor"
)

type exponentialRetryStrategy struct {
	*baseStrategy
	initialBackoffSeconds int
}

func NewExponentialRetryStrategy(maxRetries int, initialBackoffSeconds int) (exponentialRetryStrategy, error) {
	if initialBackoffSeconds <= 0 {
		return exponentialRetryStrategy{}, errors.New("no negative backoff allowed")
	}

	strategy, err := NewBaseStrategy(maxRetries)

	return exponentialRetryStrategy{baseStrategy: strategy, initialBackoffSeconds: initialBackoffSeconds}, err
}

func (strategy exponentialRetryStrategy) Apply(processRequest requestProcessor) (textProcessor.Tweets, error) {
	println("exponentialRetryStrategy")
	var resp *http.Response
	var err error

	waitTime := time.Duration(strategy.initialBackoffSeconds) * time.Second

	for strategy.currentAttempt = 1; strategy.currentAttempt <= strategy.maxRetries; strategy.currentAttempt++ {
		resp, err = processRequest(strategy.currentAttempt)

		if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			return strategy.sanitizeBodyToTweets(resp.Body), nil
		}

		time.Sleep(waitTime)
		waitTime *= 2
	}

	if err != nil {
		return nil, err
	} else {
		return nil, errors.New("failed to get a valid response")
	}
}
