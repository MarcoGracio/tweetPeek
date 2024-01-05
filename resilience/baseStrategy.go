package resilience

import (
	"errors"
	"net/http"
)

type RequestProcessor func(int) (*http.Response, error)

type IStrategy interface {
	Apply(RequestProcessor) ([]byte, error)
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
