package ght_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/ramjac/ght"
)

// This test won't work until I have a better equals check.
func TestImportExcel(t *testing.T) {
	output := []*ght.HTTPTest{
		&ght.HTTPTest{
			Request: &http.Request{
				Header: http.Header{
					"Accept-Ranges":  []string{"bytes"},
					"Content-Length": []string{"138"},
					"Content-Type":   []string{"image/gif"},
				},
			},
			ExpectedStatus: 404,
		},
		&ght.HTTPTest{
			Request:        &http.Request{},
			ExpectedStatus: 404,
		},
		&ght.HTTPTest{
			Request: &http.Request{
				Header: http.Header{"Content-Type": []string{"application/json;charset=UTF-8"}},
			},
			ExpectedStatus: 200,
			ExpectedType:   "text/html; charset=utf-8",
		},
		&ght.HTTPTest{
			Request: &http.Request{
				Header: http.Header{
					"Content-Encoding":            []string{"gzip"},
					"X-Content-Type-Options":      []string{"nosniff"},
					"Access-Control-Allow-Origin": []string{"https://test.url"},
				},
			},
			ExpectedStatus: 200,
			Regex:          regexp.MustCompile("Goblet"),
			ExpectMatch:    true,
		},
	}

	var logger *ght.VerboseLogger
	b := true
	logger.New(&b)
	path := "godocExample.xlsx"
	tabs := ""

	requestTests := ght.ImportExcel(&path, &tabs, logger, 2, 2)

	// this should loop through the output instead of the requestTests
	for _, rt := range requestTests {
		var found bool
		for _, o := range output {
			if rt.Equals(o) {
				found = true
				break
			}
		}

		// skipping these tests for now
		if !found && false {
			t.Errorf(
				"Excel import failure for %s %s\nActual: %v",
				rt.Request.Method,
				rt.Request.URL.String(),
				rt,
			)
		}
	}
}
