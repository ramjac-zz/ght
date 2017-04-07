package ght_test

import (
	"net/http"
	"reflect"
	"regexp"
	"testing"

	"github.com/ramjac/ght"
)

// This test won't work until I have a better equals check.
func TestParseCSV(t *testing.T) {
	// need more cases...
	var headers = []csvCheck{
		{"", make([]*ght.HTTPTest, 0, 0)},
		{
			"http://localhost:8080/djjff,Accept-Ranges:bytes&Content-Length:138&Content-Type:image/gif,404",
			[]*ght.HTTPTest{
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
			},
		},
		{
			"http://localhost:8080/djjff,,404,,,,http://localhost:8080,Content-Type:application/json;charset=UTF-8,200,text/html; charset=utf-8",
			[]*ght.HTTPTest{
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
			},
		},
		{
			"http://localhost:8080,content-encoding:gzip&x-content-type-options:nosniff&access-control-allow-origin:https://test.url,200,,Goblet,true",
			[]*ght.HTTPTest{
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
			},
		},
	}

	var logger *ght.VerboseLogger
	b := true
	logger.New(&b)

	for _, h := range headers {
		var ht []*ght.HTTPTest

		ht = ght.ParseCSV(&h.input, logger, 0, 0)

		// sadly always false for the test as written.
		if !reflect.DeepEqual(ht, h.output) {
			t.Errorf(
				"CSV parsing failure for input: %s\nExpected: %v\nActual: %v",
				h.input,
				h.output,
				ht,
			)
		}
	}
}

type csvCheck struct {
	input  string
	output []*ght.HTTPTest
}
