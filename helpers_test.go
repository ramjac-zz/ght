package ght_test

import "net/url"
import "regexp"

func MustParseUrl(u string) *url.URL {
	parsed, err := url.Parse(u)

	if err != nil {
		panic(err)
	}

	return parsed
}

func MustCompileRegex(input string) *regexp.Regexp {
	r, err := regexp.Compile(input)

	if err != nil {
		panic(err)
	}

	return r
}
