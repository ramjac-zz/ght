// A quick and dirty HTTP testing application for use with things like Jenkins.

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"

	"github.com/ramjac/ght"
)

func main() {
	// read flags
	retries := flag.Int("r", 5, "Number of retries for HTTP requests (defaults to 5).")
	timeElapse := flag.Int("te", 5, "Time elapse multiplier used between HTTP request retries in seconds (defaults to 5 seconds).")
	timeOut := flag.Int("to", 5000, "Time out specifies the total time in milliseconds to give each request (defaults to 5000).")
	rawCsv := flag.String("csv", "", "<url>,<headers as key1:value1&key2:value2>,<expected HTTP status code>,<expected content type>,<regex>,<bool regex should return data>")
	jsonFile := flag.String("json", "", "Path and name of the json request file.")
	excelFile := flag.String("excel", "", "Path and name of the excel file.")
	tabs := flag.String("tabs", "", "Tabs to test in the excel file.")
	parallelism := flag.Int("p", runtime.NumCPU(), "Number of requests to make concurrently (defaults to 1)")
	verbose := flag.Bool("v", false, "Prints resutls of each step. Also causes all tests to execute instead of returning after the first failure.")

	flag.Parse()
	var logger *ght.OptionalLogger
	logger.New(verbose)

	// The documentation implies this is a bad solution
	runtime.GOMAXPROCS(*parallelism)

	var r [][]*ght.HTTPTest

	switch {
	case len(*jsonFile) > 0:
		log.Fatal("JSON file support not yet implemented")
	case len(*excelFile) > 0:
		r = ght.ImportExcel(excelFile, tabs, logger, *retries, *timeElapse, *timeOut)
	case len(*rawCsv) > 0:
		r = make([][]*ght.HTTPTest, 1)
		r[0] = ght.ParseCSV(rawCsv, logger, *retries, *timeElapse, *timeOut)
	default:
		log.Fatal("An excel, JSON, or CSV input is required")
	}

	// make HTTP requests
	var wg sync.WaitGroup
	var fm sync.Mutex
	var failures int
	var failTests []string

	// Handle cancellation
	ctx := context.Background()
	// trap Ctrl+C and call cancel on the context
	c := make(chan os.Signal, *parallelism+1)
	ctx, cancel := context.WithCancel(ctx)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	// Run the requests...
	for _, v := range r {
		wg.Add(1)

		go func(v []*ght.HTTPTest) {
			for _, h := range v {
				if !h.TryRequest(ctx, cancel, logger, &wg) {
					fm.Lock()
					failures++
					failTests = append(failTests, h.Request.URL.String())
					fm.Unlock()
				}
			}
		}(v)
	}

	wg.Wait()

	// return success/failure
	logger.Println("Failing tests:", failTests)
	os.Exit(failures)
}
