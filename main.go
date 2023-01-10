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
	r.Use(handlers.Cors)

	r.Get("/", handlers.ServeFile("./static/index.html"))
	r.Get("/about", handlers.ServeFile("./static/about/index.html"))

	r.Route("/index", func(r chi.Router) {

		r.Get("/", handlers.ServeFile("./static/app/index.html"))

		r.Route("/{indexName}", func(r chi.Router) {

			r.Get("/", handlers.ServeFile("./static/app/index.html"))
			r.Get("/search", handlers.SearchField)
			r.Post("/search", handlers.SearchAdvance)

		})

	})

	handlers.FileServer(r, "/assets", http.Dir("./static/assets"))

	http.ListenAndServe(":3000", r)
}
