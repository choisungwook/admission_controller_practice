package main

import (
	"fmt"
	"net/http"
)

type requestAndResponse map[string]string

func (m requestAndResponse) ServeHTTPIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s\n", "hello", m["hello"])
}

func (m requestAndResponse) ServeHTTPSample(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s\n", "do", m["do"])
}

func main() {
	requestSample := requestAndResponse{
		"hello": "world",
		"do":    "you see me?",
	}

	mux := http.NewServeMux()

	uriPrefix := "/api"
	mux.HandleFunc(uriPrefix+"/hello", requestSample.ServeHTTPIndex)
	mux.HandleFunc(uriPrefix+"/sample", requestSample.ServeHTTPSample)

	http.ListenAndServe(":8080", mux)
}
