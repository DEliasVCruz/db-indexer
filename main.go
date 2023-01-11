package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8000, "port to set the app server to listen to")
	flag.Parse()

	if err := check.ValidPort(port); err != nil {
		fmt.Println(err.Error())
		fmt.Println("Provide a port number between 1023 and 65535")
		os.Exit(1)
	}

	ipAddr := check.GetIP()

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

	fmt.Printf("Server is running at http://%s:%d\n", ipAddr.String(), port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
