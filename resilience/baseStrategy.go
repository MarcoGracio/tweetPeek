package resilience

import (
	"errors"
	"net/http"
)

type requestProcessor func(int) (*http.Response, error)

type IStrategy interface {
	Apply(requestProcessor) ([]byte, error)
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

func (strategy *baseStrategy) GetCurrentAttempt() int {
	return strategy.currentAttempt
}
