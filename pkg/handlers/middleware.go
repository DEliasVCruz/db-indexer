package handlers

import (
	"net/http"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		headers := w.Header()

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {

			headers.Add("Vary", "Origin")
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")

			addCors(headers)
			w.WriteHeader(http.StatusOK)

		} else {

			addCors(headers)
			next.ServeHTTP(w, r)

		}

	})
}

func addCors(headers http.Header) {
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, Authorization, Content-Length, X-Requested-With")
}
