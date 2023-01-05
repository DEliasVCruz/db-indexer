package main

import (
	"net/http"

	"github.com/DEliasVCruz/db-indexer/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/index/{indexName}/search", handlers.SearchContents)

	http.ListenAndServe(":3000", r)
}
