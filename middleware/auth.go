package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jandro-es/merlin/configs"
)

func AuthValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/api/") {
			ctx := r.Context()
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		// Get the endpoint configuration based on the request path and method
		endpointConfig, ok := configs.FindConfiguration(r.Method, r.URL.Path)
		if !ok {
			http.Error(w, "Endpoint not found", http.StatusNotFound)
			return
		}
		ctx := r.Context()

		// Is Auth required?
		if endpointConfig.Auth && endpointConfig.AuthProvider != "" {
			switch endpointConfig.AuthProvider {
			case "firebase":
				token, err := validateAuthFirebase(r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				ctx = context.WithValue(ctx, "auth_token", token)
			default:
				http.Error(w, fmt.Sprintf("Invalid authentication provider '%s'", endpointConfig.AuthProvider), http.StatusInternalServerError)
				return
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateAuthFirebase(r *http.Request) (string, error) {
	firebaseAuth := configs.FirebaseClient
	authorizationToken := r.Header.Get("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	if idToken == "" {
		return "", fmt.Errorf("token not available")
	}
	//verify token
	token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		fmt.Printf("## Error validating token: %s", err)
		return "", fmt.Errorf("invalid token")
	}
	return token.UID, nil
}
