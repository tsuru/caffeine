package main

import (
	"net/http"
	"net/http/httputil"
)

func createProxy(r *http.Request, customHeaderValue string) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL = r.URL
		req.URL.Scheme = "http"
		req.URL.Host = r.Host
		req.Host = r.Host

		if customHeaderValue != "" {
			req.Header.Add("X-Caffeine", customHeaderValue)
		}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	return proxy
}
