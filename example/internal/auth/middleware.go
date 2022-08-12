package auth

import (
	"net/http"

	"ferry"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "supersecret" {
			ferry.DefaultErrorHandler(w, r, ferry.ClientError{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized",
			})
		}

		next.ServeHTTP(w, r)
	})
}
