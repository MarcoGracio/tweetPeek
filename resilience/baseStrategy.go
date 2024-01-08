package resilience

import (
	"errors"
	"net/http"
)

type RequestProcessor func(int) (*http.Response, error)

type IStrategy interface {
	Apply(RequestProcessor, chan ResponseStrategy, chan error)
}

type baseStrategy struct {
	maxRetries int
	name       string
}

type ResponseStrategy struct {
	Body           []byte
	CurrentAttempt int
	StrategyName   string
}

func newBaseStrategy(maxRetries int, name string) (*baseStrategy, error) {
	if maxRetries <= 0 {
		return nil, errors.New("max retries must be greater than 0")
	}
	return &baseStrategy{maxRetries: maxRetries, name: name}, nil
}
