package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"tweetPeek/strategy"
)

func main() {
	retryStrategy, _ := calculateBestStrategy(rand.Intn(20))
	strResp := retryStrategy.Apply(calling)
	fmt.Println(strResp)
}

func calculateBestStrategy(random int) (strategy.IStrategy, error) {
	if random > 10 {
		return strategy.NewSimpleRetryStrategy(3)
	} else {
		return strategy.NewExponentialRetryStrategy(3, 1)
	}

}

func calling(retryNr int) (*http.Response, error) {
	if retryNr <= 1 {
		return nil, nil
	}
	return http.Get("http://example.com/")
}
