package middleware

import (
	"net/http"

	"github.com/lancer2672/Dandelion_Gateway/utils"
)

const EDITOR_API_KEY = "editor_api_key"
const ADMIN_API_KEY = "admin_api_key"
const USER_API_KEY = "user_api_key"

func CheckApiKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(utils.API_KEY)
		if apiKey != "admin_api_key" && apiKey != "editor_api_key" && apiKey != "user_api_key" {
			http.Error(w, "Invalid API Key", http.StatusUnauthorized)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
