package ght

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func ParseCSV(rawCSV *string, logger *VerboseLogger, retries, timeElapse int) (r []*HTTPTest) {
	tmpClient := new(HTTPTest)

	colCount := 0

	for _, v := range strings.Split(*rawCSV, ",") {
		v = strings.TrimSpace(v)

		switch colCount {
		case 0:
			tmpClient.request = new(http.Request)
			tmpClient.retries = retries
			tmpClient.timeElapse = timeElapse

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

			AddHTTPTest(tmpClient, &r)

			tmpClient = new(HTTPTest)

			colCount = 0
			continue
		}
		colCount++
	}

	// We'll check to see if there is an unadded tmpClient so that trailing commas aren't required.
	if tmpClient.request != nil {
		AddHTTPTest(tmpClient, &r)
	}

	return r
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
