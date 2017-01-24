package ght

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// HTTPTest is a request to be tested.
type HTTPTest struct {
	request             *http.Request
	expectedStatus      int
	expectedType        string
	regex               *regexp.Regexp
	expectMatch         bool
	retries, timeElapse int
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

func AddHTTPTest(t *HTTPTest, r *[]*HTTPTest) {
	// add the tmpClient to the slice if it is valid
	// if tmpClient is valid when it has a url and expected status code
	if t.request.URL != nil &&
		t.expectedStatus > 0 {
		*r = append(*r, t)
	}
}

func (h *HTTPTest) TryRequest(logger *VerboseLogger, c chan int, wg *sync.WaitGroup) bool {
	defer wg.Done()

	// Think the for needs to contain a select or be replaced by one.
	for tries := 0; tries < h.retries; tries++ {
		select {
		case <-c:
			return true

			// need to change this to not sleep
			// basically it should have no case to enter until the proper time has elapsed or quit happens
		default:
			time.Sleep(time.Duration(h.timeElapse) * time.Duration(tries) * time.Second)

			if h.checkRequest(logger) {
				return true
			}
		}
	}

	// signal the other go routines to cancel if not in verbose mode
	if !logger.IsVerbose() {
		c <- 1
	}

	return false
}

func (h *HTTPTest) checkRequest(logger *VerboseLogger) bool {
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
