package middleware

import (
	"net/http"
	"strings"
)

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("Authorization"), " ")[1]
		if token == "" {
			http.Error(w, "Invalid token ", http.StatusUnauthorized)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
