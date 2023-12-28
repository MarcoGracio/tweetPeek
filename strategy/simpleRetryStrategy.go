package strategy

import (
	"errors"
	"io"
	"net/http"
)

type simpleRetryStrategy struct {
	baseStrategy
	maxRetries int
}

func NewSimpleRetryStrategy(maxRetries int) (*simpleRetryStrategy, error) {
	if maxRetries <= 0 {
		return nil, errors.New("max retries must be greater than 0")
	}

	return &simpleRetryStrategy{baseStrategy{}, maxRetries}, nil
}

func (strategy *simpleRetryStrategy) Apply(calling call) string {
	var resp *http.Response
	var err error
	for i := 0; i < (*strategy).maxRetries; i++ {
		resp, err = calling(i)

		if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			responseBody, _ := io.ReadAll(resp.Body)
			tweet := strategy.baseStrategy.sanitizeHtmlToTweet(string(responseBody))
			return "simpleRetryStrategy" + "\n" + tweet
		}
	}

	if err != nil {
		return err.Error()
	} else {
		return "Failed to get a valid response."
	}
}
