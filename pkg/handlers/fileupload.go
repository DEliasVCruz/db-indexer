package handlers

import (
	"archive/tar"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 500 // 500MB

func FileUpload(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "file exceed max upload size", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "application/zip" && filetype != "application/x-gzip" {

		tr := tar.NewReader(file)

		_, err := tr.Next()
		if err != nil {
			http.Error(w, "the provided file format is not allowed", http.StatusBadRequest)
			return
		}

		filetype = "application/tar"
	}

	payload, _ := json.Marshal(`{"message": "file succesfully uploaded"}`)

	log.Printf("The file uploaded is %s\n", filetype)
	w.Write(payload)

}
