package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/DEliasVCruz/db-indexer/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(handlers.Cors)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		tsFile, err := template.ParseFiles("./static/index.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "could not parse index", http.StatusInternalServerError)
			return
		}

		if err := tsFile.Execute(w, nil); err != nil {
			log.Println(err.Error())
			http.Error(w, "could not write template response", http.StatusInternalServerError)
		}

	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {

		tsFile, err := template.ParseFiles("./static/about/index.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "could not parse index", http.StatusInternalServerError)
			return
		}

		if err := tsFile.Execute(w, nil); err != nil {
			log.Println(err.Error())
			http.Error(w, "could not write template response", http.StatusInternalServerError)
		}

	})

	r.Route("/index", func(r chi.Router) {

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {

			tsFile, err := template.ParseFiles("./static/app/index.html")
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "could not parse index", http.StatusInternalServerError)
				return
			}

			if err := tsFile.Execute(w, nil); err != nil {
				log.Println(err.Error())
				http.Error(w, "could not write template response", http.StatusInternalServerError)
			}

		})

		r.Route("/{indexName}", func(r chi.Router) {

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {

				tsFile, err := template.ParseFiles("./static/app/index.html")
				if err != nil {
					log.Println(err.Error())
					http.Error(w, "could not parse index", http.StatusInternalServerError)
					return
				}

				if err := tsFile.Execute(w, nil); err != nil {
					log.Println(err.Error())
					http.Error(w, "could not write template response", http.StatusInternalServerError)
				}

			})

			r.Get("/search", handlers.SearchField)
			r.Post("/search", handlers.SearchAdvance)

		})

	})

	FileServer(r, "/assets", http.Dir("./static/assets"))

	http.ListenAndServe(":3000", r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
