package resilience

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

func NewExponentialRetryStrategy(maxRetries int, initialBackoffSeconds int) (exponentialRetryStrategy, error) {
	if initialBackoffSeconds <= 0 {
		return exponentialRetryStrategy{}, errors.New("no negative backoff allowed")
	}

	strategy, err := newBaseStrategy(maxRetries, "exponential retry")

	return exponentialRetryStrategy{baseStrategy: strategy, initialBackoffSeconds: initialBackoffSeconds}, err
}

func (strategy exponentialRetryStrategy) Apply(processRequest RequestProcessor, chanResponse chan ResponseStrategy, chanError chan error) {
	var resp *http.Response
	var err error

	for currentAttempt := 1; currentAttempt <= strategy.maxRetries; currentAttempt++ {
		resp, err = processRequest(currentAttempt)

		if resp != nil {
			defer resp.Body.Close()

			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				body, _ := io.ReadAll(resp.Body)
				chanResponse <- ResponseStrategy{Body: body, CurrentAttempt: currentAttempt, StrategyName: strategy.name}
				return
			}
		}

		waitTime := time.Duration(strategy.initialBackoffSeconds) * time.Second
		time.Sleep(waitTime)
		waitTime *= 2
	}

	if err != nil {
		chanError <- err
	} else {
		chanError <- errors.New("failed to get a valid response")
	}
}
