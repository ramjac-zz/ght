// This webserver exists to test authentication with GHT
package main

import (
	"net/http"
)

func main() {
	// jar := &myjar{}
	// jar.jar = make(map[string][]*http.Cookie)
	// client.Jar = jar

	http.HandleFunc("/BasicAuth", cookieAuth)
	http.HandleFunc("/BasicTest", cookieTest)
	http.HandleFunc("/JwtAuth", jwtAuth)
	http.HandleFunc("/JwtTest", jwtTest)
	//http.HandleFunc("/OAuth", OAuth)
	//http.HandleFunc("/OAuthTest", OAuthTest)
	http.ListenAndServe(":8080", nil)
}
