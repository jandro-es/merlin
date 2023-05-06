package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jandro-es/merlin/configs"
	"github.com/jandro-es/merlin/models"
)

// Validates that the headers of the request matches the ones specified.
func HeadersValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !strings.Contains(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		// Get the endpoint configuration based on the request path and method
		endpointConfig, ok := configs.FindConfiguration(r.Method, r.URL.Path)
		if !ok {
			http.Error(w, "Endpoint not found", http.StatusNotFound)
			return
		}
		// Validate the supplied headers agains the definition checking their required values and if they are required or not
		for key := range endpointConfig.Headers {
			requestHeader := r.Header.Get(key)
			if endpointConfig.Headers[key].Required {
				if requestHeader == "" {
					http.Error(w, fmt.Sprintf("Missing required header '%s'", key), http.StatusBadRequest)
					return
				}
				_, err := validateHeaderValue(key, requestHeader, endpointConfig.Headers[key].Validation)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				if requestHeader != "" {
					_, err := validateHeaderValue(key, requestHeader, endpointConfig.Headers[key].Validation)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
				}
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateHeaderValue(headerKey string, headerValue string, validation models.RequestHeaderValidation) (bool, error) {
	switch validation.Type {
	case "string":
		if headerValue != validation.Value {
			return false, fmt.Errorf("header %s is not valid", headerKey)
		}
	case "uuid":
		_, err := uuid.Parse(headerValue)
		if err != nil {
			return false, fmt.Errorf("header %s is not valid UUID", headerKey)
		}
	default:
		log.Fatalf("The validation type for the header %s is not supported: %s", headerKey, validation.Type)
		os.Exit(1)
	}
	return true, nil
}
