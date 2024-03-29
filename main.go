package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/lancer2672/Dandelion_Gateway/middleware"
	"github.com/lancer2672/Dandelion_Gateway/utils"
)

type Route struct {
	PathPrefix string
	BackendURL string
}

var Routes = []Route{}
var noAuthRoutes = []string{
	"/api/auth/login",
	"/api/auth/register",
}

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	Routes = []Route{
		{"/notification/*", config.NotificationServiceAddress},
		{"/movies/*", config.MovieGRPCAddress},
		{"/api/auth/login", config.MainServiceAddress},
		{"/api/auth/register", config.MainServiceAddress},
		{"/*", config.MainServiceAddress},
	}

	r := chi.NewRouter()

	for _, route := range Routes {
		var handler http.Handler

		handler = middleware.CheckApiKey(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			forwardRequest(route.BackendURL, w, r)
		}))

		if utils.StringContains(noAuthRoutes, route.PathPrefix) {
			log.Println("NoAuthRoute", route.PathPrefix)
			handler = middleware.VerifyToken(handler)
		}

		r.Handle(route.PathPrefix, handler)
	}

	http.Handle("/", r)

	http.ListenAndServe(config.GatewayAddress, nil)
	log.Println("Server started at:", config.GatewayAddress)
}

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
	// proxy.ModifyResponse = rewriteBody
}

// func rewriteBody(resp *http.Response) (err error) {
// 	b, err := io.ReadAll(resp.Body) //Read html
// 	if err != nil {
// 		return err
// 	}
// 	err = resp.Body.Close()
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("RewriteBody", b)
// 	// b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1) // replace html
// 	// body := io.NopCloser(bytes.NewReader(b))
// 	// resp.Body = body
// 	// resp.ContentLength = int64(len(b))
// 	// resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
// 	return nil
// }
