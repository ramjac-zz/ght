package ght_test

import (
	"net/url"

	"github.com/dlclark/regexp2"
)

func MustParseUrl(u string) *url.URL {
	parsed, err := url.Parse(u)

	if err != nil {
		panic(err)
	}

	return parsed
}

func MustCompileRegex(input string) *regexp2.Regexp {
	r, err := regexp2.Compile(input, regexp2.Compiled)

	if err != nil {
		panic(err)
	}

	return r
}
