package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func forwardRequest(target string, w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(target)
	log.Println("URL", url)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Host)
	proxy.ServeHTTP(w, r)
}
