package textProcessor

import (
	"fmt"
)

type threadOfTwitts []string

func (twitts threadOfTwitts) PrintTwitts() {
	for i, twitt := range twitts {
		fmt.Printf("%v - %v \n", i+1, twitt)
	}
}
