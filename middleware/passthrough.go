package middleware

import (
	"net/http"

	"github.com/jandro-es/merlin/configs"
)

// Middleware function to pass the origin headers to the response as defined.
func PassthroughHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the endpoint configuration based on the request path and method
		endpointConfig, ok := configs.FindConfiguration(r.Method, r.URL.Path)
		if !ok {
			http.Error(w, "Endpoint not found", http.StatusNotFound)
			return
		}

		// Get the request headers that need to be passthrough to the response
		for key := range endpointConfig.Headers {
			requestHeader := r.Header.Get(key)
			if endpointConfig.Headers[key].Passthrough && requestHeader != "" {
				w.Header().Set(key, requestHeader)
			}
		}
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
