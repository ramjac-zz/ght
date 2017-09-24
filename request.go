package ght

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/dlclark/regexp2"
)

// HTTPTest is a request to be tested.
type HTTPTest struct {
	Label                        string
	Request                      *http.Request
	ExpectedStatus               int
	ExpectedType                 string
	Regex                        *regexp2.Regexp
	ExpectMatch                  bool
	Retries, TimeElapse, TimeOut int
}

// Some basic pretty printing. This could use improvement.
// func (h *HTTPTest) String() string {
// 	f := `{ %s %s
// 	Expected Status: %v
// 	Expected Type: %s
// 	Regex: %s
// 	Should Regex Match: %t }`
// 	return fmt.Sprintf(
// 		f,
// 		h.Request.Method,
// 		h.Request.URL,
// 		h.ExpectedStatus,
// 		h.ExpectedType,
// 		h.Regex,
// 		h.ExpectMatch,
// 	)
// }

// formatRequest generates ascii representation of a request
func (h *HTTPTest) String() string {
	// Create return string
	var request []string

	// Add the request string
	url := fmt.Sprintf("%v %v %v", h.Request.Method, h.Request.URL, h.Request.Proto)
	request = append(request, url)

	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", h.Request.Host))

	// Loop through headers
	for name, headers := range h.Request.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if h.Request.Method == "POST" {
		h.Request.ParseForm()
		request = append(request, "\n")
		request = append(request, h.Request.Form.Encode())
	}

	// Add the expected response info
	request = append(request, fmt.Sprintf("Expected Status: %v", h.ExpectedStatus))
	request = append(request, fmt.Sprintf("ExpectedType: %v", h.ExpectedType))
	request = append(request, fmt.Sprintf("Regex: %v", h.Regex))
	request = append(request, fmt.Sprintf("Expect Match: %v", h.ExpectMatch))

	// Return the request as a string
	return strings.Join(request, "\n")
}

// AddHTTPTest appends an HTTPTest to the given slice.
func AddHTTPTest(t *HTTPTest, r *[]*HTTPTest) {
	// add the tmpClient to the slice if it is valid
	// tmpClient is valid when it has a url and expected status code
	if t.Request.URL != nil &&
		t.ExpectedStatus > 0 {
		*r = append(*r, t)
	}
}

// TryRequest will attempt an HTTP request as many times as specifie and return true if it reaches a successful response.
func (h *HTTPTest) TryRequest(ctx context.Context, cancel func(), logger *VerboseLogger, wg *sync.WaitGroup) bool {
	defer wg.Done()
	for tries := 0; tries < h.Retries; tries++ {
		select {
		case <-ctx.Done():
			return true
		case <-time.After(time.Duration(h.TimeElapse) * time.Duration(tries) * time.Second):
			if h.checkRequest(logger) {
				return true
			}
		}
	}

	// signal the other go routines to cancel if not in verbose mode
	if !logger.IsVerbose() {
		cancel()
	}

	return false
}

func (h *HTTPTest) checkRequest(logger *VerboseLogger) bool {
	client := &http.Client{
		Timeout: (time.Duration)(h.TimeOut) * time.Millisecond,
	}
	resp, err := client.Do(h.Request)

	lr, _ := httputil.DumpRequest(h.Request, true)
	logger.Printf("Test: %s", lr)

	if err == nil && resp.StatusCode == h.ExpectedStatus {
		lr, _ = httputil.DumpResponse(resp, true)
		logger.Printf("Response: %s", lr)

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

			m, err := h.Regex.MatchString(string(tmp[:]))

			if m != h.ExpectMatch || err != nil {
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

// Equals checks to see if two HTTPTests have the same field values.
func (h *HTTPTest) Equals(c *HTTPTest) bool {
	if h.Request != nil && c.Request != nil {
		if !reflect.DeepEqual(h.Request.URL, c.Request.URL) {
			return false
		}
		if !strings.EqualFold(h.Request.Method, c.Request.Method) {
			return false
		}
		//need to check headers and body... eventually
	}
	// couldn't think of a better way to do this atm
	if (h.Request == nil && c.Request != nil) || (h.Request != nil && c.Request == nil) {
		return false
	}

	if h.ExpectedStatus != c.ExpectedStatus {
		return false
	}
	if !strings.EqualFold(h.ExpectedType, c.ExpectedType) {
		return false
	}
	if h.Retries != c.Retries {
		return false
	}
	if h.TimeElapse != c.TimeElapse {
		return false
	}
	if h.ExpectMatch != c.ExpectMatch {
		return false
	}
	if h.TimeOut != c.TimeOut {
		return false
	}

	// if !strings.EqualFold(ht.Regex.String(), c.Regex.String()) {
	// 	return false
	// }

	return true
}
