package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/serviceA":
			serveReverseProxy("http://localhost:8081", w, r)
		case "/serviceB":
			serveReverseProxy("http://localhost:8082", w, r)
		case "/serviceC":
			serveReverseProxy("http://localhost:8083", w, r)
		default:
			http.Error(w, "Service not found", 404)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveReverseProxy(target string, w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	proxy.ServeHTTP(w, r)
}
