package ght_test

import "net/url"

func MustParseUrl(u string) *url.URL {
	parsed, _ := url.Parse(u)
	return parsed
}
