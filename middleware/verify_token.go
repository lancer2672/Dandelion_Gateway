package middleware

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		accessToken := strings.Split(authorizationHeader, " ")[1]
		if accessToken == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		publicKeyPem, err := os.ReadFile("public.key")
		if err != nil {
			http.Error(w, "Error when reading public key", http.StatusInternalServerError)
			return
		}

		block, _ := pem.Decode(publicKeyPem)
		if block == nil {
			http.Error(w, "Failed to parse PEM block containing the public key", http.StatusInternalServerError)
			return
		}
		publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			http.Error(w, "Failed to parse DER encoded public key", http.StatusInternalServerError)
			return
		}
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})
		if err != nil {
			http.Error(w, "Failed to verify token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Failed to get token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "payload", claims)
		r = r.WithContext(ctx)
		fmt.Println(claims)
		next.ServeHTTP(w, r)
	})
}
