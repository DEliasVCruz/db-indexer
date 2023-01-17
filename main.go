package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	var port int
	var indexDir string
	var name string
	flag.IntVar(&port, "port", 8000, "port to set the app server to listen to")
	flag.StringVar(&indexDir, "index-dir", "", "directory to index before starting server")
	flag.StringVar(&name, "index name", "", "name of your first index")
	flag.Parse()

	if err := check.ValidPort(port); err != nil {
		log.Println(err.Error())
		log.Println("Provide a port number between 1023 and 65535")
		os.Exit(1)
	}

	ipAddr := check.GetIP()

	if indexDir != "" {

		if name == "" {
			name = "MyIndex"
		}

		if err := handlers.CreateDirIndex(indexDir, name); err != nil {
			log.Println(err.Error())
		}
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(handlers.Cors)

	r.Get("/", handlers.ServeFile("./static/index.html"))
	r.Get("/about", handlers.ServeFile("./static/about/index.html"))

	r.Route("/index", func(r chi.Router) {

		r.Get("/", handlers.ServeFile("./static/app/index.html"))

		r.Get("/{indexName}", handlers.ServeFile("./static/app/index.html"))

	})

	r.Route("/api", func(r chi.Router) {

		r.Route("/index/{indexName}", func(r chi.Router) {

			r.Get("/search", handlers.SearchField)
			r.Get("/status", handlers.SearchIndexStatus)
			r.Post("/search", handlers.SearchAdvance)
			r.Put("/upload", handlers.FileUpload)

		})

	})

	handlers.FileServer(r, "/assets", http.Dir("./static/assets"))

	fmt.Printf("Server is running at http://%s:%d\n", ipAddr.String(), port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
