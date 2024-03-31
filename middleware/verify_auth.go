package middleware

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
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

func handleRefreshToken(refreshToken string, credential interface{}, clientID string) (context.Context, error) {
	payload, err := verifyToken(refreshToken)

	if err != nil || clientID != payload["userId"].(string) {
		return nil, errors.New("Invalid Request")
	}

	ctx := context.WithValue(context.Background(), "credential", credential)
	ctx = context.WithValue(ctx, "userId", clientID)
	ctx = context.WithValue(ctx, "refreshToken", refreshToken)

	return ctx, nil
}

func handleAccessToken(accessToken string, credential interface{}, clientID string) (context.Context, error) {
	payload, err := verifyToken(accessToken)

	if err != nil || clientID != payload["userId"].(string) {
		return nil, errors.New("Invalid Request")
	}

	ctx := context.WithValue(context.Background(), "credential", credential)
	ctx = context.WithValue(ctx, "userId", clientID)

	return ctx, nil
}

func VerifyAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Header.Get("X-Client-ID")
		if clientID == "" {
			http.Error(w, "Invalid Request", http.StatusUnauthorized)
			return
		}

		user, err := services.UserService.FindByID(clientID)
		if err != nil {
			http.Error(w, "Not Found User", http.StatusNotFound)
			return
		}

		credential, err := services.CredentialService.FindByID(user.Credential)
		if err != nil {
			http.Error(w, "Not Found Credential", http.StatusNotFound)
			return
		}

		refreshToken := r.Header.Get("X-Refresh-Token")
		if refreshToken != "" {
			ctx, err := handleRefreshToken(refreshToken, credential, clientID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			http.Error(w, "Invalid Request", http.StatusUnauthorized)
			return
		}

		ctx, err := handleAccessToken(accessToken, credential, clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
