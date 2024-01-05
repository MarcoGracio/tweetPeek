package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"tweetPeek/resilience"
	"tweetPeek/textProcessor"
)

func main() {
	process(requestProcessor)

	fmt.Printf("\n\n")

	process(func(currentAttempt int) (*http.Response, error) {
		return http.Get("https://en.wikipedia.org/wiki/Portugal")
	})
}

func process(request resilience.RequestProcessor) {
	retryStrategy, _ := calculateBestStrategy(rand.Intn(20))
	responseBody, err := retryStrategy.Apply(request)
	if err != nil {
		fmt.Println("fatal error")
		return
	}
	pageTwitter := textProcessor.NewPageTwitter(3, 80, responseBody)
	fmt.Printf("Attempt nr: %v \n", strconv.Itoa(retryStrategy.GetCurrentAttempt()))
	pageTwitter.ThreadOfTwitts.PrintTwitts()
}

func calculateBestStrategy(random int) (resilience.IStrategy, error) {
	if random > 10 {
		return resilience.NewSimpleRetryStrategy(3)
	} else {
		return resilience.NewExponentialRetryStrategy(3, 1)
	}
}

func requestProcessor(currentAttempt int) (*http.Response, error) {
	if currentAttempt < 2 {
		return nil, nil
	}
	return http.Get("http://example.com/")
}
