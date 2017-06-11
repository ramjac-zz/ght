package ght_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/ramjac/ght"
)

// This test won't work until I have a better equals check.
func TestImportExcel(t *testing.T) {
	// table of tests
	output := []*ght.HTTPTest{
		// GET Tests
		&ght.HTTPTest{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    MustParseUrl("http://localhost:6060"),
			},
			ExpectedStatus: 200,
			Retries:        2,
			TimeElapse:     2,
			TimeOut:        2,
		},
		&ght.HTTPTest{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    MustParseUrl("http://localhost:6060"),
				Header: http.Header{
					"accepts": {"text/html; charset=utf-8"},
				},
			},
			ExpectedStatus: 200,
			Retries:        2,
			TimeElapse:     2,
			TimeOut:        2,
		},
		&ght.HTTPTest{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    MustParseUrl("http://localhost:6060"),
				Header: http.Header{
					"accepts": {"text/html; charset=utf-8"},
				},
			},
			ExpectedStatus: 200,
			ExpectedType:   "text/html; charset=utf-8",
			Regex:          MustCompileRegex("(?i)(download go)"),
			ExpectMatch:    true,
			Retries:        2,
			TimeElapse:     2,
			TimeOut:        2,
		},
		// POST Tests
		&ght.HTTPTest{
			Request: &http.Request{
				Method:        http.MethodPost,
				URL:           MustParseUrl("http://127.0.0.1:3999/fmt"),
				Body:          ioutil.NopCloser(strings.NewReader("body=package+main%0A%0Aimport+%22fmt%22%0A%0Afunc+main()+%7B%0A%09fmt.Println(%22Hello%2C+%E4%B8%96%E7%95%8C%22)%0A%7D%0A&imports=false")),
				ContentLength: 140,
			},
			ExpectedStatus: 200,
			Retries:        2,
			TimeElapse:     2,
			TimeOut:        2,
		},
		&ght.HTTPTest{
			Request: &http.Request{
				Method: http.MethodPost,
				URL:    MustParseUrl("http://127.0.0.1:3999/fmt"),
				Header: http.Header{
					"Host":             {"127.0.0.1:3999"},
					"User-Agent":       {"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0"},
					"Accept":           {"application/json, text/plain, */*"},
					"Accept-Language":  {"en-US,en;q=0.5"},
					"Accept-Encoding":  {"gzip, deflate"},
					"Referer":          {"http://127.0.0.1:3999/welcome/1"},
					"x-requested-with": {"XMLHttpRequest"},
					"Content-Type":     {"application/x-www-form-urlencoded"},
					"Connection":       {"keep-alive"},
				},
				Body:          ioutil.NopCloser(strings.NewReader("body=package+main%0A%0Aimport+%22fmt%22%0A%0Afunc+main()+%7B%0A%09fmt.Println(%22Hello%2C+%E4%B8%96%E7%95%8C%22)%0A%7D%0A&imports=false")),
				ContentLength: 140,
			},
			ExpectedStatus: 200,
			Retries:        2,
			TimeElapse:     2,
			TimeOut:        2,
		},
		&ght.HTTPTest{
			Request: &http.Request{
				Method: http.MethodPost,
				URL:    MustParseUrl("http://127.0.0.1:3999/fmt"),
				Header: http.Header{
					"Host":             {"127.0.0.1:3999"},
					"User-Agent":       {"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0"},
					"Accept":           {"application/json, text/plain, */*"},
					"Accept-Language":  {"en-US,en;q=0.5"},
					"Accept-Encoding":  {"gzip, deflate"},
					"Referer":          {"http://127.0.0.1:3999/welcome/1"},
					"x-requested-with": {"XMLHttpRequest"},
					"Content-Type":     {"application/x-www-form-urlencoded"},
					"Connection":       {"keep-alive"},
				},
				Body:          ioutil.NopCloser(strings.NewReader("body=package+main%0A%0Aimport+%22fmt%22%0A%0Afunc+main()+%7B%0A%09fmt.Println(%22Hello%2C+%E4%B8%96%E7%95%8C%22)%0A%7D%0A&imports=false")),
				ContentLength: 140,
			},
			ExpectedStatus: 200,
			ExpectedType:   "text/plain; charset=utf-8",
			Regex:          MustCompileRegex("(fmt.Println)"),
			ExpectMatch:    true,
			Retries:        2,
			TimeElapse:     2,
			TimeOut:        2,
		},
	}

	// setup
	var logger *ght.VerboseLogger
	b := true
	logger.New(&b)
	path := "godocExample.xlsx"
	tabs := ""

	requestTests := ght.ImportExcel(&path, &tabs, logger, 2, 2, 2)

	// run tests
	for _, o := range output {
		var found bool
		for _, rt := range requestTests {
			if rt.Equals(o) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf(
				"Excel import failure for %s",
				o,
			)
		}
	}
}
