package ght

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// HTTPTest is a request to be tested.
type HTTPTest struct {
	Request             *http.Request
	ExpectedStatus      int
	ExpectedType        string
	Regex               *regexp.Regexp
	ExpectMatch         bool
	Retries, TimeElapse int
}

func (h *HTTPTest) String() string {
	f := `%s %s
	Expected Status: %v
	Expected Type: %s
	Regex: %s
	Should Regex Match: %t`
	return fmt.Sprintf(
		f,
		h.Request.Method,
		h.Request.URL,
		h.ExpectedStatus,
		h.ExpectedType,
		h.Regex,
		h.ExpectMatch,
	)
}

// AddHTTPTest appends an HTTPTest to the given slice.
func AddHTTPTest(t *HTTPTest, r *[]*HTTPTest) {
	// add the tmpClient to the slice if it is valid
	// if tmpClient is valid when it has a url and expected status code
	if t.Request.URL != nil &&
		t.ExpectedStatus > 0 {
		*r = append(*r, t)
	}
}

// TryRequest will attempt an HTTP request as many times as specified and return true if it reaches a successful response.
func (h *HTTPTest) TryRequest(logger *VerboseLogger, c chan int, wg *sync.WaitGroup) bool {
	defer wg.Done()

	// Think the for needs to contain a select or be replaced by one.
	for tries := 0; tries < h.Retries; tries++ {
		select {
		case <-c:
			return true

			// need to change this to not sleep
			// basically it should have no case to enter until the proper time has elapsed or quit happens
		default:
			time.Sleep(time.Duration(h.TimeElapse) * time.Duration(tries) * time.Second)

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
	resp, err := client.Do(h.Request)

	logger.Printf("Test - %v", h)

	if err == nil &&
		resp.StatusCode == h.ExpectedStatus {
		logger.Printf("Response: %v", *resp)

		if len(h.ExpectedType) > 0 &&
			!strings.EqualFold(resp.Header.Get("content-type"), h.ExpectedType) {
			return false
		}

		if h.Regex != nil {
			tmp, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				logger.Printf("Body could not be read: %v", err)
				return false
			}
			//Should I be reading the response body differently here...

			m := h.Regex.MatchString(string(tmp[:]))

			if m != h.ExpectMatch {
				logger.Printf("Body: \"%s\"", string(tmp))

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
