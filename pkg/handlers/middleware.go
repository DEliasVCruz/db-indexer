package handlers

import (
	"log"
	"net/http"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {

			handlePreflight(w, r)
			w.WriteHeader(http.StatusOK)

		} else {

			handleNormalRequest(w, r)
			next.ServeHTTP(w, r)

		}

	})
}

func addCors(headers http.Header) {
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, Authorization, Content-Length, X-Requested-With")
}

func handleNormalRequest(w http.ResponseWriter, r *http.Request) {

	headers := w.Header()
	origin := r.Header.Get("Origin")

	headers.Add("Vary", "Origin")
	if origin == "" {
		log.Println("No provided origin")
		return
	}

	addCors(headers)

}

func handlePreflight(w http.ResponseWriter, r *http.Request) {

	headers := w.Header()
	origin := r.Header.Get("Origin")

	if origin == "" {
		log.Println("No provided origin")
		return
	}

	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")

	addCors(headers)

}
