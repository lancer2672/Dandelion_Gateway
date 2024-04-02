package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/lancer2672/Dandelion_Gateway/internal/middleware"
	"github.com/lancer2672/Dandelion_Gateway/internal/utils"
)

type Route struct {
	PathPrefix string
	BackendURL string
}

var Routes = []Route{}
var noAuthRoutes = []string{
	"/api/auth/login",
	"/api/auth/register",
	//TODO: handle auth
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestLimitter)
	for _, route := range Routes {
		var handler http.Handler

		handler = (http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			forwardRequest(route.BackendURL, w, r)
		}))

		if utils.StringContains(noAuthRoutes, route.PathPrefix) {
			log.Println("NoAuthRoute", route.PathPrefix)
			handler = (middleware.VerifyAuthentication(handler))
		}

		r.Handle(route.PathPrefix, handler)
	}
	return r
}
