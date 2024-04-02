package middleware

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lancer2672/Dandelion_Gateway/internal/api"
	"github.com/pkg/errors"
)

func verifyToken(accessToken string) (jwt.MapClaims, error) {
	publicKeyPem, err := os.ReadFile("public.key")
	if err != nil {

		return nil, errors.New("Error when reading public key")
	}

	block, _ := pem.Decode(publicKeyPem)
	if block == nil {
		return nil, errors.New("Failed to parse PEM block containing the public key")
	}
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Failed to parse DER encoded public key")
	}
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, errors.New("Failed to verify token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Failed to get token claims")
	}
	return claims, nil
}

func handleVerifyRefreshToken(refreshToken string, clientID string) error {
	payload, err := verifyToken(refreshToken)

	if err != nil || clientID != payload["userId"].(string) {
		return errors.New("Invalid Request")
	}
	return nil
}

func handleVerifyAccessToken(accessToken string, clientID string) error {
	payload, err := verifyToken(accessToken)

	if err != nil || clientID != payload["userId"].(string) {
		return errors.New("Invalid Request")
	}

	return nil
}

func VerifyAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("x-client-id")
		if userId == "" {
			http.Error(w, "Invalid Request", http.StatusUnauthorized)
			return
		}
		//check user credential is exist
		_, err := api.GetUserCredential(userId)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
		}

		refreshToken := r.Header.Get("x-refresh-token")
		//if this request is token refresh request then let it pass
		if refreshToken != "" {
			err := handleVerifyRefreshToken(refreshToken, userId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			http.Error(w, "Invalid Request", http.StatusUnauthorized)
			return
		}

		err = handleVerifyAccessToken(accessToken, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
