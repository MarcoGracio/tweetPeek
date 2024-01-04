package textProcessor

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type PageTwitter struct {
	nrMaxtwitts    int
	twittLength    int
	ThreadOfTwitts threadOfTwitts
}

func NewPageTwitter(nrMaxtwitts int, twittLenght int, data []byte) PageTwitter {
	textCleaned := PageTwitter{}.sanitizeByteSlice(data)
	threadOfTwitts := PageTwitter{}.getThreadOfTwitts(textCleaned, nrMaxtwitts, twittLenght)
	return PageTwitter{nrMaxtwitts: nrMaxtwitts, twittLength: twittLenght, ThreadOfTwitts: threadOfTwitts}
}

func (pageTwitter PageTwitter) sanitizeByteSlice(data []byte) string {
	p := bluemonday.StrictPolicy()
	htmlSanitized := p.Sanitize(string(data))

	fields := strings.Fields(string(htmlSanitized))
	return strings.Join(fields, " ")
}

func (PageTwitter) getThreadOfTwitts(htmlText string, nrMax int, twittLength int) threadOfTwitts {
	var twitt string
	var threadOfTwitts threadOfTwitts
	for i := 0; i < len(htmlText); i += twittLength {

		if (i + twittLength) > len(htmlText) {
			twitt = htmlText[i:]
		} else {
			twitt = htmlText[i:(i + twittLength)]
		}

		threadOfTwitts = append(threadOfTwitts, twitt)
	}

	if nrMax <= len(threadOfTwitts) {
		return threadOfTwitts[:nrMax]
	}
	return threadOfTwitts
}
