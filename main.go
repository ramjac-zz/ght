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
	"runtime"
	"strconv"
	"strings"
	"sync"
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

func (h *HTTPTest) String() string {
	f := `URL: %s
	Expected Status: %v
	Expected Type: %s
	Regex: %s
	Should Regex Match: %t`
	return fmt.Sprintf(
		f,
		h.request.URL,
		h.expectedStatus,
		h.expectedType,
		h.regex,
		h.expectMatch,
	)
}

var (
	verbose                       bool
	retries, timeElapse, failures int
	logger                        *verboseLogger
)

func main() {
	logger = new(verboseLogger)

	// read flags
	flag.IntVar(&retries, "r", 5, "Number of retries for HTTP requests (defaults to 5).")
	flag.IntVar(&timeElapse, "t", 5, "Time elapse multiplier used between HTTP request retries in seconds (defaults to 5).")
	rawCsv := flag.String("csv", "", "<url>,<headers as key1:value1&key2:value2>,<expected HTTP status code>,<expected content type>,<regex>,<bool regex should return data>")
	jsonFile := flag.String("json", "", "Path and name of the json request file.")
	concurrency := flag.Int("c", 1, "Number of requests to make concurrently (defaults to 1)")
	flag.BoolVar(&verbose, "v", false, "Prints resutls of each step. Also causes all tests to execute instead of returning after the first failure.")

	flag.Parse()

	// The documentation implies this is a bad solution
	runtime.GOMAXPROCS(*concurrency)

	var r []*HTTPTest

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
	var wg sync.WaitGroup
	var fm sync.Mutex
	c := make(chan int, *concurrency+1)

	// Run the requests...
	for _, v := range r {
		wg.Add(1)

		go v.tryRequest(c, &fm, &wg)
	}

	wg.Wait()

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
				logger.Println(err)
			}
		case 1:
			tmpClient.setHeaders(v)
		case 2:
			s, err := strconv.Atoi(v)
			if err == nil {
				tmpClient.expectedStatus = s
			} else {
				logger.Printf("Error parsing status code: %s\n", err)
			}
		case 3:
			tmpClient.expectedType = v
		case 4:
			if len(v) > 0 {
				s, err := regexp.Compile(v)
				if err != nil {
					logger.Printf("Error parsing regular expression: %s\n", err)
				} else {
					tmpClient.regex = s
				}
			}
		case 5:
			if len(v) > 0 {
				s, err := strconv.ParseBool(v)
				if err != nil {
					logger.Printf("Error parsing the boolean for whether the regex should match or not: %s\n", err)
				} else {
					tmpClient.expectMatch = s
				}
			}

			addHTTPTest(tmpClient, &r)

			tmpClient = new(HTTPTest)

			colCount = 0
			continue
		}
		colCount++
	}

	// We'll check to see if there is an unadded tmpClient so that trailing commas aren't required.
	if tmpClient.request != nil {
		addHTTPTest(tmpClient, &r)
	}

	return r
}

func addHTTPTest(t *HTTPTest, r *[]*HTTPTest) {
	// add the tmpClient to the slice if it is valid
	// if tmpClient is valid when it has a url and expected status code
	if t.request.URL != nil &&
		t.expectedStatus > 0 {
		*r = append(*r, t)
	}
}

func (h *HTTPTest) tryRequest(quit chan int, fm *sync.Mutex, wg *sync.WaitGroup) {
	// Think the for needs to contain a select or be replaced by one.
	for tries := 0; tries < retries; tries++ {
		select {
		case <-quit:
			wg.Done()
			return

			// need to change this to not sleep
			// basically it should have no case to enter until the proper time has elapsed or quit happens
		default:
			time.Sleep(time.Duration(timeElapse) * time.Duration(tries) * time.Second)

			if h.checkRequest() {
				wg.Done()
				return
			}
		}
	}

	fm.Lock()
	failures++
	fm.Unlock()

	// break on the first failure if not in verbose mode
	if !verbose {
		quit <- 1
	}

	wg.Done()
}

func (h *HTTPTest) checkRequest() bool {
	client := &http.Client{}
	resp, err := client.Do(h.request)

	logger.Printf("Test - %v", h)

	if err == nil &&
		resp.StatusCode == h.expectedStatus {
		logger.Printf("Response: %v", *resp)

		if len(h.expectedType) > 0 &&
			strings.Compare(resp.Header.Get("content-type"), h.expectedType) != 0 {
			return false
		}

		if h.regex != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)

			m := h.regex.MatchString(buf.String())

			if m != h.expectMatch {
				return false
			}
		}
		return true
	}

	if err != nil {
		logger.Printf("Error on Get: %s\n", err)
	} else {
		logger.Printf("Error in response: %v\n", *resp)
	}
	return false
}

func (h *HTTPTest) setHeaders(headerString string) {
	headers := strings.Split(headerString, "&")
	h.request.Header = make(map[string][]string)
	for _, tmp := range headers {
		kv := strings.SplitN(tmp, ":", 2)

		if len(kv) != 2 {
			continue
		}

		h.request.Header.Set(kv[0], kv[1])
	}
}

// verboseLogger only logs when the verbose variable is true.
type verboseLogger struct{}

func (l *verboseLogger) Println(v ...interface{}) {
	if verbose {
		log.Println(v)
	}
}

func (l *verboseLogger) Printf(s string, v ...interface{}) {
	if verbose {
		log.Printf(s, v)
	}
}
