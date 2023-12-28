package strategy

import (
	"errors"
	"io"
	"net/http"
	"time"
)

type exponentialRetryStrategy struct {
	baseStrategy
	maxRetries            int
	initialBackoffSeconds int
}

func NewExponentialRetryStrategy(maxRetries int, initialBackoffSeconds int) (*exponentialRetryStrategy, error) {
	if maxRetries <= 0 || initialBackoffSeconds <= 0 {
		return nil, errors.New("incorrect arguments")
	}

	return &exponentialRetryStrategy{baseStrategy: baseStrategy{}, maxRetries: maxRetries, initialBackoffSeconds: initialBackoffSeconds}, nil
}

func (strategy *exponentialRetryStrategy) Apply(calling call) string {
	var resp *http.Response
	var err error

	waitTime := time.Duration(strategy.initialBackoffSeconds) * time.Second

	for i := 0; i < strategy.maxRetries; i++ {
		resp, err = calling(i)

		if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			responseBody, _ := io.ReadAll(resp.Body)
			tweet := strategy.baseStrategy.sanitizeHtmlToTweet(string(responseBody))
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
