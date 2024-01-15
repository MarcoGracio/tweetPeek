package textProcessor

import "testing"

var testsNewPageTwitter = []struct {
	name        string
	nrMaxTweets int
	tweetLenght int
	data        []byte
	expected    string
	isErr       bool
}{
	{"newPageTwitter-validData", 1, 7, []byte("<p>valid    data valid data valid data valid data</p>"), "valid d", false},
	{"newPageTwitter-inValidData", -1, -1, nil, "", true},
}

func Test_newPageTwitter(t *testing.T) {
	for _, tt := range testsNewPageTwitter {
		got, err := NewPageTwitter(tt.nrMaxTweets, tt.tweetLenght, tt.data)

		if tt.isErr {
			if err == nil {
				t.Errorf("%s expected an error but did not get one", tt.name)
			}
		} else {
			if err != nil {
				t.Errorf("%s did not expect an error but get one", err.Error())
			}
		}

		for _, tweet := range got.ThreadOfTweets {
			if tweet != tt.expected {
				t.Errorf("expected '%v' but got '%v'", tt.expected, tweet)
			}
		}
	}
}

var testSanitizeByteSlice = []struct {
	name     string
	data     []byte
	expected string
	isErr    bool
}{
	{"sanitizeByteSlice-removeTags", []byte("<h1>Hello<p> </p> World</h1>"), "Hello World", false},
	{"sanitizeByteSlice-removeBlanks", []byte("Hello      World"), "Hello World", false},
	{"sanitizeByteSlice-invalidData-Nil", nil, "", true},
	{"sanitizeByteSlice-invalidData-Empty", []byte(""), "", true},
}

func Test_sanitizeByteSlice(t *testing.T) {
	for _, tt := range testSanitizeByteSlice {
		stringSanitized, err := PageTwitter{}.sanitizeByteSlice(tt.data)

		if tt.isErr {
			if err == nil {
				t.Errorf("%s expected an error but did not get one", tt.name)
			}
		} else {
			if err != nil {
				t.Errorf("%s did not expect an error but get one", err.Error())
			}
		}

		if stringSanitized != tt.expected {
			t.Errorf("expected '%v' but got '%v'", tt.expected, stringSanitized)
		}
	}
}

var testGetThreadOfTweets = []struct {
	name        string
	text        string
	nrMaxTweets int
	tweetLenght int
	expected    []string
	isErr       bool
}{
	{"getThreadOfTweets-validText", "this is a valid text", 2, 4, []string{"this", " is "}, false},
	{"getThreadOfTweets-smallerTextThanTweetlenght", "this is a valid text", 2, 100, []string{"this is a valid text"}, false},
	{"getThreadOfTweets-smallerTextThanNrmaxtweets", "this", 2, 3, []string{"thi", "s"}, false},
	{"getThreadOfTweets-invalidNrMaxTwitts", "", -1, 20, []string{""}, true},
	{"getThreadOfTweets-invalidTwettLenght", "", 1, -20, []string{""}, true},
}

func Test_getThreadOfTweets(t *testing.T) {
	for _, tt := range testGetThreadOfTweets {
		tweets, err := PageTwitter{}.getThreadOfTweets(tt.text, tt.nrMaxTweets, tt.tweetLenght)

		if tt.isErr {
			if err == nil {
				t.Errorf("%s expected an error but did not get one", tt.name)
			}
		} else {
			if err != nil {
				t.Errorf("%s did not expect an error but get one", err.Error())
			}
		}

		for i, tweet := range tweets {
			if tweet != tt.expected[i] {
				t.Errorf("expected '%v' but got '%v'", tt.expected[i], tweet)
			}
		}
	}
}
