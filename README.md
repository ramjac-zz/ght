Go HTTP test aims to make it easy to create automated HTTP test scripts.

The csv looks like this:
"<url>,<headers as key1:value1&key2:value2>,<expected HTTP status code>,<expected content type>,<regex>,<bool should regex match>"

Some examples run with godoc -http=:8080

    go run main.go -r 1 -t 1 -csv "http://localhost:8080/djjff,,404,,,,http://localhost:8080,,200,text/html; charset=utf-8,,,http://localhost:8080,,200,,(Download go),true"

Case insentive example
    go run main.go -v -r 1 -t 1 -csv "http://localhost:8080,,200,,(?i)(download go),true"


A nice little reference for Regex as parsed by Golang
https://regex-golang.appspot.com/assets/html/index.html

 TODO

 * Implement concurrent requests.
 * Allow for a JSON file input of whose schema is based on an array of the http.Request struct. I'd like for there to be more flexibility in creating the requests.
 * Improve testing for CSV inputs
 * Improve verbose output

