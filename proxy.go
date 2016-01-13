package main

import (
	"net/http"
	"net/http/httputil"
)

func createProxy(r *http.Request) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL = r.URL
		req.URL.Scheme = "http"
		req.URL.Host = r.Host
		req.Host = r.Host
	}
	proxy := &httputil.ReverseProxy{Director: director}
	return proxy
}
