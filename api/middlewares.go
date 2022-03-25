package api

import (
	"net/http"
)


func BodyLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		}
		next.ServeHTTP(w, r)
	})
}