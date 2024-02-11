package main

import (
	"fmt"
	"net/http"
)

type requestAndResponse map[string]string

func (m requestAndResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range m {
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}
}

func main() {
	requestSample := requestAndResponse{
		"hello": "world",
		"do":    "you see me?",
	}
	http.ListenAndServe(":8080", requestSample)
}
