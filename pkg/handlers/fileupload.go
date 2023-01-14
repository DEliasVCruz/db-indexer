package handlers

import (
	"archive/tar"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/index"
)

const MAX_UPLOAD_SIZE = 500<<20 + 1024 // 500MB
const MAX_MEMORY = 20 << 20            // 20MB

func FileUpload(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		http.Error(w, "file exceed max upload size", http.StatusBadRequest)
		return
	}

	file, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		file.Close()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		file.Close()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "application/zip" && filetype != "application/x-gzip" {

		tr := tar.NewReader(file)

		_, err := tr.Next()
		if err != nil {
			file.Close()
			http.Error(w, "the provided file format is not allowed", http.StatusBadRequest)
			return
		}

		filetype = "application/tar"
	}

	go index.NewIndex(
		"My Index",
		strings.TrimPrefix(filetype, "application/"),
		"",
		&data.FormFile{File: file, Size: fh.Size},
	)

	response, err := json.Marshal(&FileUploaded{Uploaded: true, State: "processing"})
	if err != nil {
		http.Error(w, "server marshaling failed", http.StatusInternalServerError)
		return
	}

	w.Write(response)

}
