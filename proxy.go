package main

import (
	"net/http"
	"net/http/httputil"
)

func createProxy(r *http.Request, customHeaderValue string) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL = r.URL
		req.URL.Scheme = "http"
		host := r.Header.Get("X-Host")
		req.URL.Host = host
		req.Host = host

		if customHeaderValue != "" {
			req.Header.Add("X-Caffeine", customHeaderValue)
		}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	return proxy
}
