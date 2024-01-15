package textProcessor

import (
	"fmt"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type PageTwitter struct {
	ThreadOfTweets threadOfTweets
}

func NewPageTwitter(nrMaxTweets int, tweetLength int, data []byte) (PageTwitter, error) {
	textCleaned, errSanitize := PageTwitter{}.sanitizeByteSlice(data)
	if errSanitize != nil {
		return PageTwitter{}, errSanitize
	}
	threadOfTweets, errThreadOfTweets := PageTwitter{}.getThreadOfTweets(textCleaned, nrMaxTweets, tweetLength)
	if errThreadOfTweets != nil {
		return PageTwitter{}, errThreadOfTweets
	}
	return PageTwitter{ThreadOfTweets: threadOfTweets}, nil
}

func (PageTwitter) sanitizeByteSlice(data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("incorrect arguments: data=%v", data)
	}

	p := bluemonday.StrictPolicy()
	htmlSanitized := p.Sanitize(string(data))

	fields := strings.Fields(string(htmlSanitized))
	return strings.Join(fields, " "), nil
}

func (PageTwitter) getThreadOfTweets(text string, nrMaxTweets int, tweetLength int) (threadOfTweets, error) {
	if nrMaxTweets <= 0 || tweetLength <= 0 {
		return threadOfTweets{}, fmt.Errorf("incorrect arguments: nrMaxTweets=%d, tweetLength=%d", nrMaxTweets, tweetLength)
	}

	var tweet string
	var threadOfTweets threadOfTweets
	for i := 0; i < len(text) && len(threadOfTweets) < nrMaxTweets; i += tweetLength {

		if (i + tweetLength) > len(text) {
			tweet = text[i:]
		} else {
			tweet = text[i:(i + tweetLength)]
		}

		threadOfTweets = append(threadOfTweets, tweet)
	}

	return threadOfTweets, nil
}
