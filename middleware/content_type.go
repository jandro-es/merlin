package middleware

import "net/http"

func ContentTypeApplicationJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
