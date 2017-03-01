// A quick and dirty HTTP testing application for use with things like Jenkins.

package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/ramjac/ght"
)

func main() {
	// read flags
	retries := flag.Int("r", 5, "Number of retries for HTTP requests (defaults to 5).")
	timeElapse := flag.Int("t", 5, "Time elapse multiplier used between HTTP request retries in seconds (defaults to 5).")
	rawCsv := flag.String("csv", "", "<url>,<headers as key1:value1&key2:value2>,<expected HTTP status code>,<expected content type>,<regex>,<bool regex should return data>")
	jsonFile := flag.String("json", "", "Path and name of the json request file.")
	excelFile := flag.String("excel", "", "Path and name of the excel file.")
	tabs := flag.String("tabs", "", "Tabs to test in the excel file.")
	parallelism := flag.Int("p", runtime.NumCPU(), "Number of requests to make concurrently (defaults to 1)")
	verbose := flag.Bool("v", false, "Prints resutls of each step. Also causes all tests to execute instead of returning after the first failure.")

	flag.Parse()
	var logger *ght.VerboseLogger
	logger.New(verbose)

	// The documentation implies this is a bad solution
	runtime.GOMAXPROCS(*parallelism)

	var r []*ght.HTTPTest

	switch {
	case len(*jsonFile) > 0:
		log.Fatal("JSON file support not yet implemented")
	case len(*excelFile) > 0:
		r = ght.ImportExcel(excelFile, tabs, logger, *retries, *timeElapse)
	case len(*rawCsv) > 0:
		// parse csv to structs
		r = ght.ParseCSV(rawCsv, logger, *retries, *timeElapse)
	default:
		log.Fatal("A JSON or CSV input is required")
	}

	// make HTTP requests
	var wg sync.WaitGroup
	var fm sync.Mutex
	var failures int
	c := make(chan int, *parallelism+1)
	var failTests []string

	// Run the requests...
	for _, v := range r {
		wg.Add(1)

		go func(v *ght.HTTPTest) {
			if !v.TryRequest(logger, c, &wg) {
				fm.Lock()
				failures++
				failTests = append(failTests, v.Request.URL.String())
				fm.Unlock()
			}
		}(v)
	}

	wg.Wait()

	// return success/failure
	logger.Println("Failing tests:", failTests)
	fmt.Println(failures)
}
