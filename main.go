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

func main() {
	// read flags
	retries := flag.Int("r", 5, "Number of retries for HTTP requests (defaults to 5)")
	timeElapse := flag.Int("t", 5, "Time elapse multiplier used between HTTP request retries in seconds (defaults to 5)")
	rawCsv := flag.String("csv", "", "<url>,<headers as key1=value1&key2=value2>,<expected HTTP status code>,<expected content type>,<regex>,<bool regex should return data>")
	//concurrency := flag.Int("c", 0, "Number of requests to make concurrently (defaults to 1)")

	flag.Parse()

	// read csv

	// parse csv to structs
	var r []*HTTPTest
	tmpClient := new(HTTPTest)

	colCount := 0

	for _, v := range strings.Split(*rawCsv, ",") {
		v = strings.TrimSpace(v)

		switch colCount {
		case 0:
			tmpClient.request = new(http.Request)
			u, err := url.Parse(v)
			if err == nil {
				tmpClient.request.URL = u
			} else {
				log.Println(err)
			}
		case 1:
			setHeaders(tmpClient, v)
		case 2:
			s, err := strconv.Atoi(v)
			if err == nil {
				tmpClient.expectedStatus = s
			} else {
				log.Printf("Error parsing status code: %s\n", err)
			}
		case 3:
			tmpClient.expectedType = v
		case 4:
			if len(v) > 0 {
				s, err := regexp.Compile(v)
				if err != nil {
					log.Printf("Error parsing regular expression: %s\n", err)
				} else {
					tmpClient.regex = s
				}
			}
		case 5:
			if len(v) > 0 {
				s, err := strconv.ParseBool(v)
				if err != nil {
					log.Printf("Error parsing the boolean for whether the regex should match or not: %s\n", err)
				} else {
					tmpClient.expectMatch = s
				}
			}
			// add the tmpClient to the slice if it is valid
			// if tmpClient has a url and expected status code set then append to r

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
	}

	// return success/failure
	fmt.Println(failures)
}

func checkRequest(ht *HTTPTest) bool {
	//httptest.
	// This could/shoudl be rewritten to use http/httptest package
	//fmt.Println(ht.request)
	client := &http.Client{}
	//fmt.Println(ht.request)
	resp, err := client.Do(ht.request)
	log.Printf("Response: %v", resp)
	if err == nil &&
		resp.StatusCode == ht.expectedStatus &&
		(strings.Compare(resp.Header.Get("content-type"), ht.expectedType) == 0 || len(ht.expectedType) == 0) {

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
		log.Printf("Error on Get: %s\n", err)
	} else {
		log.Printf("Error in response: %v\n", resp)
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
			log.Println(err)
		}
	}
}
