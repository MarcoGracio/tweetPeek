package resilience

import (
	"errors"
	"fmt"
	"net/http"
	"tweetPeek/textProcessor"
)

type simpleRetryStrategy struct {
	*baseStrategy
}

func NewSimpleRetryStrategy(maxRetries int) (simpleRetryStrategy, error) {
	strategy, err := NewBaseStrategy(maxRetries)
	return simpleRetryStrategy{baseStrategy: strategy}, err
}

func (strategy simpleRetryStrategy) Apply(processRequest requestProcessor) (textProcessor.Tweets, error) {
	fmt.Println("simpleRetryStrategy")
	var resp *http.Response
	var err error
	for strategy.currentAttempt = 1; strategy.currentAttempt <= strategy.maxRetries; strategy.currentAttempt++ {
		resp, err = processRequest(strategy.currentAttempt)

		if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			return strategy.sanitizeBodyToTweets(resp.Body), nil
		}
	}

	if err != nil {
		return nil, err
	} else {
		return nil, errors.New("failed to get a valid response")
	}
}
