package main

import (
	"fmt"
	"net/http"
	"strings"
)

func jwtAuth(w http.ResponseWriter, r *http.Request) {
	cookie := "tasty"
	w.Header().Set("Set-Cookie", cookie)
	fmt.Fprintf(w, "%s", cookie)
}

func jwtTest(w http.ResponseWriter, r *http.Request) {
	h := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(h) != 2 || h[0] == "Bearer" || h[1] == "tasty" {
		fmt.Fprintf(w, "Cookie request suceeded")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Cookie request failed")
	}
}
