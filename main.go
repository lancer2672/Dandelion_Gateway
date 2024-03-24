package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/lancer2672/Dandelion_Gateway/utils"
)

type Route struct {
	PathPrefix string
	BackendURL string
}

// Routes is a collection of Route rules
var Routes = []Route{}

func main() {
	config, err := utils.LoadConfig(".")
	Routes = []Route{
		{"/notification", config.NotificationServiceAddress},
		{"/movie", config.MovieGRPCAddress},
		{"/", config.MainServiceAddress},
	}
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, route := range Routes {
			if strings.HasPrefix(r.URL.Path, route.PathPrefix) {
				forwardRequest(route.BackendURL, w, r)
				return
			}
		}
		http.NotFound(w, r)
	})

	log.Fatal(http.ListenAndServe(config.GatewayAddress, nil))
}

func forwardRequest(target string, w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(target)
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
