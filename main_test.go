package main

import (
	"net/http"
	"reflect"
	"testing"
)

func TestSetHeaders(t *testing.T) {
	// need more cases...
	var headers = []headerCheck{
		{"", http.Header{}},
	}

	for _, h := range headers {
		var ht *HTTPTest

		setHeaders(ht, h.input)

		if !reflect.DeepEqual(ht.request.Header, h.output) {
			t.Errorf("Header mismatch for input: %s", h.input)
		}
	}
}

type headerCheck struct {
	input  string
	output http.Header
}

func TestParseCSV(t *testing.T) {
	// need more cases...
	var headers = []csvCheck{
		{"", []*HTTPTest{}},
	}

	for _, h := range headers {
		var ht []*HTTPTest

		ht = parseCSV(&h.input)

		if !reflect.DeepEqual(ht, h.output) {
			t.Errorf("CSV parsing failure for input: %s", h.input)
		}
	}
}

type csvCheck struct {
	input  string
	output []*HTTPTest
}

func TestcheckRequest(t *testing.T) {
}
