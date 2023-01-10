package handlers

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"
)

func ServeFile(file string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tsFile, err := template.ParseFiles(file)

		if err != nil {

			log.Println(err.Error())
			http.Error(w, "could not parse index", http.StatusInternalServerError)
			return

		}

		if err := tsFile.Execute(w, nil); err != nil {

			log.Println(err.Error())
			http.Error(w, "could not write template response", http.StatusInternalServerError)

		}

	}

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
