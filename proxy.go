package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func createProxy(r *http.Request) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL = r.URL
		req.URL.Scheme = "http"
		req.URL.Host = r.Host
		req.Host = r.Host
		log.Printf("Request: %#v", req)
	}
	proxy := &httputil.ReverseProxy{Director: director}
	return proxy
}
