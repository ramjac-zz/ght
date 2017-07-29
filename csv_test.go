package ght_test

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/dlclark/regexp2"
	"github.com/ramjac/ght"
)

// This test won't work until I have a better equals check.
func TestParseCSV(t *testing.T) {
	// table of tests
	var tests = []csvCheck{
		{"", make([]*ght.HTTPTest, 0, 0)},
		{
			"http://localhost:8080/fail404,Accept-Ranges:bytes&Content-Length:138&Content-Type:image/gif,404",
			[]*ght.HTTPTest{
				&ght.HTTPTest{
					Request: &http.Request{
						Method: http.MethodGet,
						URL:    MustParseUrl("http://localhost:8080/fail404"),
						Header: http.Header{
							"Accept-Ranges":  {"bytes"},
							"Content-Length": {"138"},
							"Content-Type":   {"image/gif"},
						},
					},
					ExpectedStatus: 404,
					Retries:        1,
					TimeElapse:     1,
					TimeOut:        1,
				},
			},
		},
		{
			"http://localhost:8080/test2,,404,,,,http://localhost:8080,Content-Type:application/json;charset=UTF-8,200,text/html; charset=utf-8",
			[]*ght.HTTPTest{
				&ght.HTTPTest{
					Request: &http.Request{
						Method: http.MethodGet,
						URL:    MustParseUrl("http://localhost:8080/test2"),
					},
					ExpectedStatus: 404,
					Retries:        1,
					TimeElapse:     1,
					TimeOut:        1,
				},
				&ght.HTTPTest{
					Request: &http.Request{
						Method: http.MethodGet,
						URL:    MustParseUrl("http://localhost:8080"),
						Header: http.Header{"Content-Type": []string{"application/json;charset=UTF-8"}},
					},
					ExpectedStatus: 200,
					ExpectedType:   "text/html; charset=utf-8",
					Retries:        1,
					TimeElapse:     1,
					TimeOut:        1,
				},
			},
		},
		{
			"http://localhost:8080,content-encoding:gzip&x-content-type-options:nosniff&access-control-allow-origin:https://test.url,200,,Goblet,true",
			[]*ght.HTTPTest{
				&ght.HTTPTest{
					Request: &http.Request{
						Method: http.MethodGet,
						URL:    MustParseUrl("http://localhost:8080"),
						Header: http.Header{
							"Content-Encoding":            []string{"gzip"},
							"X-Content-Type-Options":      []string{"nosniff"},
							"Access-Control-Allow-Origin": []string{"https://test.url"},
						},
					},
					ExpectedStatus: 200,
					Regex:          regexp2.MustCompile("Goblet", regexp2.Compiled),
					ExpectMatch:    true,
					Retries:        1,
					TimeElapse:     1,
					TimeOut:        1,
				},
			},
		},
	}

	// setup
	logger := log.New(os.Stdout, "GHT: ", log.Lshortfile)
	// b := true
	// logger.New(&b)``

	// run tests
	for _, test := range tests {
		t.Logf("Searching for test output %v", test)

		var results []*ght.HTTPTest

		results = ght.ParseCSV(&test.input, logger, 1, 1, 1)

		for _, result := range results {
			var found bool
			for _, o := range test.output {
				if result.Equals(o) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf(
					"CSV parsing failure for input: %s\nExpected: %v\nActual: %v",
					test.input,
					test.output,
					results,
				)
			}
		}
	}
}

type csvCheck struct {
	input  string
	output []*ght.HTTPTest
}
