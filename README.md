

"<url>,<headers as key1=value1&key2=value2>,<expected HTTP status code>,<expected content type>,<regex>,<bool should regex match>"


 go run main.go -r 1 -t 1 -csv "http://localhost:8080/djjff,,404,,,,http://localhost:8080,,200,text/html; charset=utf-8,,,http://localhost:8080,,200,,[Go project],true"