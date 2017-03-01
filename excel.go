package ght

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"strings"

	"github.com/Luxurioust/excelize"
)

// ImportExcel takes an excel of the correct format and returns a slice of HTTPTest.
func ImportExcel(fileName, tabsToTest *string, logger *VerboseLogger, retries, timeElapse int) (r []*HTTPTest) {
	xlsx, err := excelize.OpenFile(*fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tabs := xlsx.GetSheetMap()

	for _, tab := range tabs {
		// here is where we could check to see that the specified tab is one that was listed.

		if len(*tabsToTest) > 0 {
			if !strings.Contains(*tabsToTest, tab) {
				continue
			}
		}

		// Get all the rows in a sheet.
		rows := xlsx.GetRows(tab)
		for _, row := range rows {
			tmpClient := new(HTTPTest)
			// range over the cells
			for k, v := range row {
				switch k {
				case 1:
					tmpClient.Request = new(http.Request)

					// need to move this to new columns
					tmpClient.Retries = retries
					tmpClient.TimeElapse = timeElapse

					u, err := url.Parse(v)
					if err == nil {
						tmpClient.Request.URL = u
					} else {
						logger.Println(err)
					}
				case 2:
					tmpClient.setHeaders(v)
				case 3:
					tmpClient.Request.Method = v
				case 4:
					tmpClient.Request.Body.Read([]byte(v))
				case 5:
					s, err := strconv.Atoi(v)
					if err == nil {
						tmpClient.ExpectedStatus = s
					} else {
						logger.Printf("Error parsing status code: %s\n", err)
					}
				case 6:
					tmpClient.ExpectedType = v
				case 7:
					if len(v) > 0 {
						s, err := regexp.Compile(v)
						if err != nil {
							logger.Printf("Error parsing regular expression: %s\n", err)
						} else {
							tmpClient.Regex = s
						}
					}
				case 8:
					if len(v) > 0 {
						s, err := strconv.ParseBool(v)
						if err != nil {
							logger.Printf("Error parsing the boolean for whether the regex should match or not: %s\n", err)
						} else {
							tmpClient.ExpectMatch = s
						}
					}
				case 9:
					s, err := strconv.Atoi(v)
					if err == nil {
						tmpClient.Retries = s
					} else {
						logger.Printf("Error parsing retries: %s\n", err)
					}
				case 10:
					s, err := strconv.Atoi(v)
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
