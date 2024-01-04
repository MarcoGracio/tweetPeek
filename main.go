package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"tweetPeek/resilience"
)

func main() {
	retryStrategy, _ := calculateBestStrategy(rand.Intn(20))
	tweetsRes, err := retryStrategy.Apply(requestProcessor)
	if err != nil {
		fmt.Println("fatal error")
		return
	}
	fmt.Printf("Attempt nr: %v \n%v", strconv.Itoa(retryStrategy.GetCurrentAttempt()), strings.Join(tweetsRes, "\n"))

	retryStrategy, _ = calculateBestStrategy(rand.Intn(20))
	tweetsRes, err = retryStrategy.Apply(func(currentAttempt int) (*http.Response, error) {
		return http.Get("https://en.wikipedia.org/wiki/Portugal")
	})
	if err != nil {
		fmt.Println("fatal error")
		return
	}
	fmt.Printf("Attempt nr: %v \n%v", strconv.Itoa(retryStrategy.GetCurrentAttempt()), strings.Join(tweetsRes, "\n"))
}

func calculateBestStrategy(random int) (resilience.IStrategy, error) {
	return resilience.NewSimpleRetryStrategy(3)
	// if random > 10 {
	// 	return strategy.NewSimpleRetryStrategy(3)
	// } else {
	// 	return strategy.NewExponentialRetryStrategy(3, 1)
	// }
}

func requestProcessor(currentAttempt int) (*http.Response, error) {
	if currentAttempt < 2 {
		return nil, nil
	}
	return http.Get("http://example.com/")
}
