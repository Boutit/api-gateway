package middleware

import (
	"net/http"
)

type Middleware = func(http.Handler) http.Handler

func CreateHandler(middleware []Middleware, mux http.Handler) http.Handler {
	combined := mux
	for i := len(middleware) - 1; i >=0; i-- {
		handler := middleware[i]
		combined = handler(combined)
	}
	return combined
}