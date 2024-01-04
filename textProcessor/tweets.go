package textProcessor

type Tweets []string

func (tweets Tweets) Get(htmlText string, nrMax int) Tweets {
	const tweetLength = 80
	var tweet string
	for i := 0; i < len(htmlText); i += tweetLength {

		if (i + tweetLength) > len(htmlText) {
			tweet = htmlText[i:]
		} else {
			tweet = htmlText[i:(i + tweetLength)]
		}

		tweets = append(tweets, tweet)
	}

	if nrMax <= len(tweets) {
		return tweets[:nrMax]
	}
	return tweets
}
