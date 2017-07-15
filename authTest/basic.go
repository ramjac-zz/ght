package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func cookieAuth(w http.ResponseWriter, r *http.Request) {
	cookie := "tasty"
	w.Header().Set("Set-Cookie", cookie)
	fmt.Fprintf(w, "Set Cookie reponse: %s", cookie)
}

func cookieTest(w http.ResponseWriter, r *http.Request) {
	h := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(h) != 2 || h[0] == "Basic" || h[1] == "tasty" {
		fmt.Fprintf(w, "Cookie request suceeded")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Cookie request failed")
	}
}

type myjar struct {
	jar map[string][]*http.Cookie
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	fmt.Printf("The URL is : %s\n", u.String())
	fmt.Printf("The cookie being set is : %s\n", cookies)
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	fmt.Printf("The URL is : %s\n", u.String())
	fmt.Printf("Cookie being returned is : %s\n", p.jar[u.Host])
	return p.jar[u.Host]
}
