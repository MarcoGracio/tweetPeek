package resilience

import (
	"errors"
	"io"
	"net/http"
)

type simpleRetryStrategy struct {
	*baseStrategy
}

func NewSimpleRetryStrategy(maxRetries int) (simpleRetryStrategy, error) {
	strategy, err := newBaseStrategy(maxRetries, "simple retry")
	return simpleRetryStrategy{baseStrategy: strategy}, err
}

func (strategy simpleRetryStrategy) Apply(processRequest RequestProcessor, chanResponse chan ResponseStrategy, chanError chan error) {
	var resp *http.Response
	var err error
	for currentAttempt := 1; currentAttempt <= strategy.maxRetries; currentAttempt++ {
		resp, err = processRequest(currentAttempt)

		if resp != nil {
			defer resp.Body.Close()

			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				body, _ := io.ReadAll(resp.Body)
				chanResponse <- ResponseStrategy{Body: body, CurrentAttempt: currentAttempt, StrategyName: "simple retry"}
				return
			}
		}
	}

	if err != nil {
		chanError <- err
	} else {
		chanError <- errors.New("failed to get a valid response")
	}
}
