package textProcessor

import (
	"fmt"
)

type threadOfTweets []string

func (tweets threadOfTweets) PrintTweets() {
	for i, tweet := range tweets {
		fmt.Printf("%v - %v \n", i+1, tweet)
	}
}
