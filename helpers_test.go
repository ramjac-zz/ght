package ght_test

import "net/url"

func MustParseUrl(u string) *url.URL {
	parsed, err := url.Parse(u)

	if err != nil {
		panic(err)
	}

	return parsed
}
