package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

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
		{"/v1/create_movie_history", config.MovieGRPCAddress},
		{"/api/auth/send-email-verification", config.MainServiceAddress},
	}
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("MAIN PATH", r.URL.Path)

	// 	switch {
	// 	case strings.HasPrefix(r.URL.Path, "/notification"):
	// 		forwardRequest(config.NotificationServiceAddress, w, r)
	// 	case strings.HasPrefix(r.URL.Path, "/movie"):
	// 		forwardRequest(config.MovieGRPCAddress, w, r)
	// 	default:
	// 		forwardRequest(config.MainServiceAddress, w, r)
	// 	}
	// })

	for _, route := range Routes {
		http.Handle(route.PathPrefix, newProxy(route.BackendURL))
	}
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

func newProxy(targetURL string) http.Handler {
	fmt.Println("TargetURL", targetURL)
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal("Error parsing target URL:", err)
	}

	return httputil.NewSingleHostReverseProxy(target)
}
