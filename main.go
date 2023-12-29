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
	fmt.Println("Attempt nr:" + strconv.Itoa(retryStrategy.GetCurrentAttempt()))
	strResp := retryStrategy.Apply(requestProcessor)
	fmt.Println("Attempt nr:" + strconv.Itoa(retryStrategy.GetCurrentAttempt()))
	fmt.Println(strResp)
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
