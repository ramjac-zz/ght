package main

import (
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

func TestSetHeaders(t *testing.T) {
	// need more cases...
	var headers = []headerCheck{
		{"", http.Header{}},
		{
			"Content-Type:application/json;charset=UTF-8",
			http.Header{"Content-Type": []string{"application/json;charset=UTF-8"}},
		},
		{
			"Accept-Ranges:bytes&Content-Length:138&Content-Type:image/gif",
			http.Header{
				"Accept-Ranges":  []string{"bytes"},
				"Content-Length": []string{"138"},
				"Content-Type":   []string{"image/gif"},
			},
		},
		{
			"content-encoding:gzip&x-content-type-options:nosniff&access-control-allow-origin:https://test.url",
			http.Header{
				"Content-Encoding":            []string{"gzip"},
				"X-Content-Type-Options":      []string{"nosniff"},
				"Access-Control-Allow-Origin": []string{"https://test.url"},
			},
		},
	}

	for _, h := range headers {
		ht := new(HTTPTest)
		ht.request = new(http.Request)

		setHeaders(ht, h.input)

		if !reflect.DeepEqual(ht.request.Header, h.output) {
			t.Errorf(
				"Header mismatch for input: %s\nExpected %s\nGot %s",
				h.input,
				h.output,
				ht.request.Header,
			)
		}
	}
}

type headerCheck struct {
	input  string
	output http.Header
}

// This test won't work until I have a better equals check.
func TestParseCSV(t *testing.T) {
	// need more cases...
	var headers = []csvCheck{
		{"", make([]*HTTPTest, 0, 0)},
		{
			"http://localhost:8080/djjff,,404",
			[]*HTTPTest{
				&HTTPTest{
					request:        &http.Request{},
					expectedStatus: 404,
				},
			},
		},
		{
			"http://localhost:8080/djjff,,404,,,,http://localhost:8080,,200,text/html; charset=utf-8",
			[]*HTTPTest{
				&HTTPTest{
					request:        &http.Request{},
					expectedStatus: 404,
				},
				&HTTPTest{
					request:        &http.Request{},
					expectedStatus: 200,
					expectedType:   "text/html; charset=utf-8",
				},
			},
		},
		{
			"http://localhost:8080,,200,,Goblet,true",
			[]*HTTPTest{
				&HTTPTest{
					request:        &http.Request{},
					expectedStatus: 200,
					regex:          regexp.MustCompile("Goblet"),
					expectMatch:    true,
				},
			},
		},
	}

	for _, h := range headers {
		var ht []*HTTPTest

		ht = parseCSV(&h.input)

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
	output []*HTTPTest
}

func TestcheckRequest(t *testing.T) {
}
