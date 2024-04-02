package middleware

import (
	"net/http"
	"strconv"

	"github.com/lancer2672/Dandelion_Gateway/internal/utils"
	"github.com/lancer2672/Dandelion_Gateway/services"
)

func RequestLimitter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userIP := r.Header.Get("X-Forwarded-For")
		if userIP == "" {
			userIP = r.RemoteAddr
		}
		value, err := services.GetValue(userIP)
		if err != nil {
			//first time query
			value = "1"
		}
		requestCount, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		if requestCount > utils.ConfigIns.RequestLimitPerTimeUnit {
			services.SetValue(userIP, strconv.Itoa(requestCount+1), utils.ConfigIns.RequestLimitTimeUnit)
			next.ServeHTTP(w, r)

		}
	})
}
