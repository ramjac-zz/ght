// A quick and dirty HTTP testing application for use with things like Jenkins.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HTTPTest is a request to be tested.
type HTTPTest struct {
	request        *http.Request
	expectedStatus int
	expectedType   string
	regex          *regexp.Regexp
	expectMatch    bool
}

var (
	verbose bool
	info    *maybeLogger
)

func main() {
	// read flags
	retries := flag.Int("r", 5, "Number of retries for HTTP requests (defaults to 5).")
	timeElapse := flag.Int("t", 5, "Time elapse multiplier used between HTTP request retries in seconds (defaults to 5).")
	rawCsv := flag.String("csv", "", "<url>,<headers as key1=value1&key2=value2>,<expected HTTP status code>,<expected content type>,<regex>,<bool regex should return data>")
	jsonFile := flag.String("json", "", "Path and name of the json request file.")
	//concurrency := flag.Int("c", 0, "Number of requests to make concurrently (defaults to 1)")
	flag.BoolVar(&verbose, "v", false, "Prints resutls of each step. Also causes all tests to execute instead of returning after the first failure.")

	flag.Parse()

	var r []*HTTPTest
	info = new(maybeLogger)

	switch {
	case len(*jsonFile) > 0:
		log.Fatal("JSON file support not yet implemented")
	case len(*rawCsv) > 0:
		// parse csv to structs
		r = parseCSV(rawCsv)
	default:
		log.Fatal("A JSON or CSV input is required")
	}

	// make HTTP requests
	failures := 0

Tests:
	for _, v := range r {
		for tries := 0; tries < *retries; tries++ {
			time.Sleep(time.Duration(*timeElapse) * time.Duration(tries) * time.Second)

			if checkRequest(v) {
				continue Tests
			}
		}
		failures++
		if !verbose {
			break
		}
	}

	// return success/failure
	fmt.Println(failures)
}

func parseCSV(rawCSV *string) (r []*HTTPTest) {
	tmpClient := new(HTTPTest)

	colCount := 0

	for _, v := range strings.Split(*rawCSV, ",") {
		v = strings.TrimSpace(v)

		switch colCount {
		case 0:
			tmpClient.request = new(http.Request)
			u, err := url.Parse(v)
			if err == nil {
				tmpClient.request.URL = u
			} else {
				info.Println(err)
			}
		case 1:
			setHeaders(tmpClient, v)
		case 2:
			s, err := strconv.Atoi(v)
			if err == nil {
				tmpClient.expectedStatus = s
			} else {
				info.Printf("Error parsing status code: %s\n", err)
			}
		case 3:
			tmpClient.expectedType = v
		case 4:
			if len(v) > 0 {
				s, err := regexp.Compile(v)
				if err != nil {
					info.Printf("Error parsing regular expression: %s\n", err)
				} else {
					tmpClient.regex = s
				}
			}
		case 5:
			if len(v) > 0 {
				s, err := strconv.ParseBool(v)
				if err != nil {
					info.Printf("Error parsing the boolean for whether the regex should match or not: %s\n", err)
				} else {
					tmpClient.expectMatch = s
				}
			}
			// add the tmpClient to the slice if it is valid
			// if tmpClient is valid when it has a url and expected status code
			if tmpClient.request.URL != nil &&
				tmpClient.expectedStatus > 0 {
				r = append(r, tmpClient)
			}

			tmpClient = new(HTTPTest)

			colCount = 0
			continue
		}
		colCount++
	}

	return r
}

func checkRequest(ht *HTTPTest) bool {
	//httptest.
	// This could/shoudl be rewritten to use http/httptest package
	//fmt.Println(ht.request)
	client := &http.Client{}
	//fmt.Println(ht.request)
	resp, err := client.Do(ht.request)
	info.Printf("Response: %v", *resp)
	if err == nil &&
		resp.StatusCode == ht.expectedStatus {
		if len(ht.expectedType) > 0 &&
			strings.Compare(resp.Header.Get("content-type"), ht.expectedType) != 0 {
			return false
		}

		if ht.regex != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)

			m := ht.regex.MatchString(buf.String())

			if m != ht.expectMatch {
				return false
			}
		}
		return true
	}

	if err != nil {
		info.Printf("Error on Get: %s\n", err)
	} else {
		info.Printf("Error in response: %v\n", resp)
	}
	return false
}

func setHeaders(ht *HTTPTest, h string) {
	// Expects a string h like "accepts=application/json&bearer=blahblahblah"
	// Values should be urlEncoded

	headers := strings.Split(h, "&")
	ht.request.Header = make(map[string][]string)
	for _, tmp := range headers {
		kv := strings.Split(tmp, "=")

		if len(kv) != 2 {
			continue
		}

		// need to fix
		k, err := url.QueryUnescape(kv[0])
		v, err := url.QueryUnescape(kv[1])

		if err == nil {
			ht.request.Header.Set(k, v)
		} else {
			info.Println(err)
		}
	}
}

// maybeLogger only logs when the verbose variable is true. I don't like the name, but I'm tired.
type maybeLogger struct{}

func (l *maybeLogger) Println(v ...interface{}) {
	if verbose {
		log.Println(v)
	}
}

func (l *maybeLogger) Printf(s string, v ...interface{}) {
	if verbose {
		log.Printf(s, v)
	}
}
