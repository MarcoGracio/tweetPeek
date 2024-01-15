// todo:
// 2- tests
// 3- containers
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
	requestProcessors := []resilience.RequestProcessor{
		requestExampleProcessor,
		requestExampleProcessor,
		func(currentAttempt int) (*http.Response, error) {
			return http.Get("https://en.wikipedia.org/wiki/Portugal")
		},
	}

	process(requestProcessors)
}

func process(requestProcessors []resilience.RequestProcessor) {

	chanResponse := make([]chan resilience.ResponseStrategy, len(requestProcessors))
	chanErrors := make([]chan error, len(requestProcessors))

	for i, processors := range requestProcessors {
		retryStrategy, _ := calculateBestStrategy()

		chanResponse[i] = make(chan resilience.ResponseStrategy)
		chanErrors[i] = make(chan error)

		go retryStrategy.Apply(processors, chanResponse[i], chanErrors[i])

	}

	for i := range requestProcessors {
		select {
		case err := <-chanErrors[i]:
			if err != nil {
				fmt.Println("fatal error")
				return
			}
		case responseStrategy := <-chanResponse[i]:
			pageTwitter, err := textProcessor.NewPageTwitter(2, 80, responseStrategy.Body)
			if err != nil {
				fmt.Println("fatal error")
				return
			}
			fmt.Printf("Strategy: %v - Attempt nr: %v \n", responseStrategy.StrategyName, strconv.Itoa(responseStrategy.CurrentAttempt))
			pageTwitter.ThreadOfTweets.PrintTweets()
		}
	}
}

func calculateBestStrategy() (resilience.IStrategy, error) {
	if rand.Intn(20) > 10 {
		return resilience.NewSimpleRetryStrategy(3)
	} else {
		return resilience.NewExponentialRetryStrategy(3, 1)
	}
}

func requestExampleProcessor(currentAttempt int) (*http.Response, error) {
	if currentAttempt < 2 {
		return nil, nil
	}
	return http.Get("http://example.com/")
}
