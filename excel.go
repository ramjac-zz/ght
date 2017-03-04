package ght

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

// ImportExcel takes an excel of the correct format and returns a slice of HTTPTest.
func ImportExcel(fileName, tabsToTest *string, logger *VerboseLogger, retries, timeElapse int) (r []*HTTPTest) {
	xlFile, err := xlsx.OpenFile(*fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, tab := range xlFile.Sheets {
		// here is where we could check to see that the specified tab is one that was listed.
		if len(*tabsToTest) > 0 {
			logger.Println("test tabs")
			if !strings.Contains(*tabsToTest, tab.Name) {
				continue
			}
		}

		for _, row := range tab.Rows {
			tmpClient := new(HTTPTest)
			// range over the cells
			for k, v := range row.Cells {
				if v == nil || strings.TrimSpace(v.Value) == "" {
					if k == 1 {
						break
					}
					continue
				}

				switch k {
				case 1:
					tmpClient.Request = new(http.Request)

					// need to move this to new columns
					tmpClient.Retries = retries
					tmpClient.TimeElapse = timeElapse

					u, err := url.Parse(v.Value)
					if err == nil {
						tmpClient.Request.URL = u
					} else {
						logger.Println(err)
					}
				case 2:
					tmpClient.setExcelHeaders(v.Value)
				case 3:
					tmpClient.Request.Method = v.Value
				case 4:
					tmpClient.Request.Body = ioutil.NopCloser(strings.NewReader(v.Value))
					tmpClient.Request.ContentLength = int64(len(v.Value))
				case 5:
					s, err := strconv.Atoi(v.Value)
					if err == nil {
						tmpClient.ExpectedStatus = s
					} else {
						logger.Printf("Error parsing status code: %s\n", err)
					}
				case 6:
					tmpClient.ExpectedType = strings.TrimSpace(v.Value)
				case 7:
					if len(v.Value) > 0 {
						s, err := regexp.Compile(v.Value)
						if err != nil {
							logger.Printf("Error parsing regular expression: %s\n", err)
						} else {
							tmpClient.Regex = s
						}
					}
				case 8:
					if len(v.Value) > 0 {
						s, err := strconv.ParseBool(v.Value)
						if err != nil {
							logger.Printf("Error parsing the boolean for whether the regex should match or not: %s\n", err)
						} else {
							tmpClient.ExpectMatch = s
						}
					}
				case 9:
					s, err := strconv.Atoi(v.Value)
					if err == nil {
						tmpClient.Retries = s
					} else {
						logger.Printf("Error parsing retries: %s\n", err)
					}
				case 10:
					s, err := strconv.Atoi(v.Value)
					if err == nil {
						tmpClient.TimeElapse = s
					} else {
						logger.Printf("Error parsing time elapse: %s\n", err)
					}

					AddHTTPTest(tmpClient, &r)

					tmpClient = new(HTTPTest)

					continue
				}
			}
			fmt.Println()
		}
	}
	return r
}

func (h *HTTPTest) setExcelHeaders(headerString string) {
	headers := strings.Split(headerString, "\n")
	h.Request.Header = make(map[string][]string)
	for _, tmp := range headers {
		kv := strings.SplitN(tmp, ":", 2)

		if len(kv) != 2 {
			continue
		}

		h.Request.Header.Set(kv[0], strings.TrimSpace(kv[1]))
	}
}
