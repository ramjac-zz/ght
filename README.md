# GHT

GHT (GHT HTTP Tester) aims to make it easy to create automated HTTP test scripts.

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Travis](https://img.shields.io/travis/rust-lang/rust.svg?style=flat-square)](https://travis-ci.org/ramjac/ght)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/ramjac/ght)
[![Go Report Card](https://goreportcard.com/badge/github.com/ramjac/ght?style=flat-square)](https://goreportcard.com/report/github.com/ramjac/ght)

## Installation:

    go get github.com/ramjac/ght/...

## Using the tool

Running an excel file looks simply like this:

    ght -excel testfile.xlsx

An example test file is provided in this repo. The test file works on local godoc and gotour servers. A flag -v will print verbose output.

The csv test format looks like this:

    <url>,<headers as key1:value1&key2:value2>,<expected HTTP status code>,<expected content type>,<regex>,<bool should regex match>

Some examples run with godoc -http=:8080

    ght -r 1 -t 1 -csv "http://localhost:8080/djjff,,404,,,,http://localhost:8080,,200,text/html; charset=utf-8,,,http://localhost:8080,,200,,(Download go),true"

Case insentive example

    ght -v -r 1 -te 1 -to 1 -csv "http://localhost:8080,,200,,(?i)(download go),true"


A nice little reference for Regex as parsed by Golang
https://regex-golang.appspot.com/assets/html/index.html

## TODO

<<<<<<< HEAD
* Allow for a JSON file input of whose schema is based on an array of the HTTPTest struct
* Improve unit tests and examples
    * Tests are currently a little redundant and only cover GET and POST
    * Request testing relies on godoc and gotour running which always breaks on Travis CI
* Improve verbose output
    * The verbose output should summarize what failed. There is a summary, but this could be more helpful.
    * The HTTP request and response should pretty print
* Fix a minor bug where spreadsheet rows that lack retries/time elapse fail to run
* Add configurable timeout to the test runner
    * Implement the timeout in the excel tester
    * Implement the timeout in the csv tester

=======
* Remove time.Sleep from TryRequest
* Allow for a JSON file input of whose schema is based on an array of the HTTPTest struct
* Improve verbose output
    * The verbose output should summarize what failed. There is a summary, but this could be more helpful.
    * The HTTP request and response should pretty print
* Fix a minor bug where spreadsheet rows that lack retries/time elapse/timeout fail to run
* Add some kind of authentication flow
    * Allow for "Set-Cookie" in a response to set the Cookies of future requests
    * Use the "token: ..." in response to set the "Authorization: " header of future requests
    * Do these two scenarios cover enough? (not cover everything, just enough)
    * How does this fit into the current test runner?
        * Perhaps CSV tests should execute serially?
        * Perhaps serially per tab?
>>>>>>> dev
