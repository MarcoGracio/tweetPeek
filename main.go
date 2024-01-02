package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"tweetPeek/strategy"
)

func main() {
	retryStrategy, _ := calculateBestStrategy(rand.Intn(20))
	strResp := retryStrategy.Apply(requestProcessor)
	fmt.Printf("Attempt nr: %v \n%v", strconv.Itoa(retryStrategy.GetCurrentAttempt()), strResp)

	retryStrategy, _ = calculateBestStrategy(rand.Intn(20))
	strResp = retryStrategy.Apply(func(currentAttempt int) (*http.Response, error) {
		return http.Get("http://google.com/")
	})
	fmt.Printf("Attempt nr: %v \n%v", strconv.Itoa(retryStrategy.GetCurrentAttempt()), strResp)

}

func calculateBestStrategy(random int) (strategy.IStrategy, error) {
	if random > 10 {
		return strategy.NewSimpleRetryStrategy(3)
	} else {
		return strategy.NewExponentialRetryStrategy(3, 1)
	}
}

func requestProcessor(currentAttempt int) (*http.Response, error) {
	if currentAttempt < 2 {
		return nil, nil
	}
	return http.Get("http://example.com/")
}
